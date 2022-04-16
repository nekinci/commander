package osutil

import (
	"commander/src/constant"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Copy struct {
	Src     string
	Dst     string
	Options CopyOptions
}

type CopyOptions struct {
	Recursive bool
	CreateDir bool // If target directory does not exist, create it
	Overwrite bool // If target file exists, overwrite it
	Depth     int  // If recursive, maximum depth
}

func CopyAll(src, dst string, recursive, createDir bool) error {

	srcStat, srcErr := os.Lstat(src)
	if srcErr != nil {
		return srcErr
	}

	dstStat, dstErr := os.Lstat(dst)
	if dstErr != nil {

		if !createDir {
			return dstErr
		}

		if srcStat.IsDir() {
			if os.IsNotExist(dstErr) {
				err := os.MkdirAll(dst, constant.DefaultDirPermission)
				if err != nil {
					return err
				}
				dstStat, err = os.Lstat(dst)
				if err != nil {
					return err
				}
			} else {
				return dstErr
			}
		}

	}

	if srcStat.IsDir() && !dstStat.IsDir() {
		return constant.NewInvalidDestError(dst, "Source is directory but destination is not. Cannot copy.")
	}

	return copyAll(src, dst, recursive)
}

func copyAll(src, dst string, recursive bool) error {
	src, err := replaceEnvironmentVariables(src)

	if err != nil {
		return err
	}

	src = addCurrentDirIfRequired(src)

	dst, err = replaceEnvironmentVariables(dst)

	if err != nil {
		return err
	}

	dst = addCurrentDirIfRequired(dst)

	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {

			commonPath := getCommonPath(src, path)
			dstPath := filepath.Join(dst, commonPath)
			dstPath = addCurrentDirIfRequired(dstPath)

			if !recursive && strings.Count(path, string(os.PathSeparator)) > 0 && src != path {
				return fs.SkipDir
			}

			return os.MkdirAll(dstPath, constant.DefaultDirPermission)

		} else {

			dir := filepath.Dir(dst)
			dir = addCurrentDirIfRequired(dir)
			err := os.MkdirAll(dir, constant.DefaultDirPermission)
			if err != nil {
				return err
			}
		}

		commonPath := getCommonPath(src, path)

		srcPath := filepath.Join(src, commonPath)
		dstPath := filepath.Join(dst, commonPath)

		return copy(srcPath, dstPath)

	})
	return err
}

func copy(source, destination string) error {
	source = addCurrentDirIfRequired(source)
	file, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	destination = addCurrentDirIfRequired(destination)
	err = os.WriteFile(destination, file, constant.DefaultFilePermission)
	if err != nil {
		return err
	}
	return nil
}

func replaceEnvironmentVariables(path string) (string, error) {
	elements := strings.Split(path, string(os.PathSeparator))
	newElements := make([]string, 0)
	for _, element := range elements {
		if strings.HasPrefix(element, "$") {
			envVar := strings.TrimPrefix(element, "$")
			envVarValue := os.Getenv(envVar)
			if envVarValue == "" {
				return "", constant.NewEnvVarNotFoundError(envVar)
			}
			newElements = append(newElements, envVarValue)
		} else if strings.HasPrefix(element, "~") {
			dir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			newElements = append(newElements, dir)
		} else if strings.HasPrefix(element, "%") && strings.HasSuffix(element, "%") {
			envVar := strings.TrimPrefix(element, "%")
			envVar = strings.TrimSuffix(envVar, "%")
			envVarValue := os.Getenv(envVar)
			if envVarValue == "" {
				return "", constant.NewEnvVarNotFoundError(envVar)
			}
			newElements = append(newElements, envVarValue)
		} else {
			newElements = append(newElements, element)
		}

	}

	return filepath.Join(newElements...), nil
}

func getCommonPath(src, path string) string {
	srcReplace := strings.Replace(src, "./", "", 1)
	res := strings.Replace(path, srcReplace, "", 1)
	return res
}

func addCurrentDirIfRequired(src string) string {
	if !strings.HasPrefix(src, "/") && !strings.HasPrefix(src, "./") {
		src = "./" + src
	}
	return src
}
