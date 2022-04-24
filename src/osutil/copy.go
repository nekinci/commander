package osutil

import (
	"bytes"
	"commander/src/constant"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FilterFunc func(fileName string) bool

type Copy struct {
	Source      string
	Destination string
	Options     *CopyOptions
	sourceStat  os.FileInfo
	destStat    os.FileInfo
}

type CopyOptions struct {
	Recursive     bool
	Depth         *int         // nil means infinite depth
	FileMode      *os.FileMode // nil means use default
	DirectoryMode *os.FileMode // nil means use default
	FilterFunc    FilterFunc
}

func NewCopy(source, destination string, copyOptions *CopyOptions) *Copy {
	return &Copy{
		Source:      source,
		Destination: destination,
		Options:     copyOptions,
	}
}

func (c *Copy) Copy() error {
	err := c.loadSourceStat()
	if err != nil {
		return err
	}

	return c.copy()
}

func (c *Copy) copy() error {

	return filepath.WalkDir(c.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err := c.validate(path); err != nil {
			return err
		}

		if d.IsDir() {

			destination := c.GetCalculatedDestination(path)
			d := c.Options.getDepthLimit()
			if d > getDepth(c.GetSource(), getPath(path)) || c.GetSource() == getPath(path) {
				err := c.mkdirHelper(destination)
				if err != nil {
					return err
				}
			}
			return nil
		}

		destination := c.GetCalculatedDestination(path)
		if strings.HasSuffix(destination, string(os.PathSeparator)) {
			if c.isDestinationDir() {
				err := c.mkdirHelper(destination)
				if err != nil {
					return err
				}
				destination += d.Name()
			} else {
				destination = strings.TrimSuffix(destination, string(os.PathSeparator))
				dir := filepath.Dir(destination)
				err := os.MkdirAll(dir, c.Options.getDirectoryMode())
				if err != nil {
					return err
				}
			}

		}

		source := getPath(path)
		return copyFile(source, destination)
	})

}

// Only internal use
func (c *Copy) mkdirHelper(destination string) error {
	lstat, err := os.Lstat(destination)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(destination, c.Options.getDirectoryMode())
			if err != nil {
				return err
			}

		}
		return err
	}

	if !lstat.IsDir() {
		return os.ErrInvalid
	}
	return nil
}

func copyFile(src, dst string) error {
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, bytes.NewReader(file))
	if err != nil {
		return err
	}

	return dstFile.Close()

}

func (c *Copy) loadSourceStat() error {
	stat, err := os.Lstat(c.Source)
	c.sourceStat = stat
	if err != nil {
		return err
	}
	return nil
}

// GetSource returns the source path as absolute path
func (c *Copy) GetSource() string {
	return getPath(c.Source)
}

// GetDestination returns the destination path as absolute path
func (c *Copy) GetDestination() string {
	return getPath(c.Destination)
}

func (c *Copy) GetCalculatedDestination(path string) string {
	dest := c.GetDestination()
	if strings.HasSuffix(dest, string(os.PathSeparator)) {
		dest = dest + c.getCommonPath(path)
	} else {
		dest = dest + string(os.PathSeparator) + c.getCommonPath(path)
	}

	return dest
}

func (c *CopyOptions) getDepthLimit() int {
	d := 0
	if c.Recursive {
		if c.Depth != nil {
			d = *c.Depth
		} else {
			d = -1
		}
	}

	return d
}

func (c *Copy) validate(path string) error {

	d := c.Options.getDepthLimit()

	if d == -1 {
		return nil // there is no depth limit
	} else if getDepth(c.GetSource(), getPath(path)) <= d {
		return nil // the depth is within the limit
	}

	return fs.SkipDir

}

func getPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	wd, _ := os.Getwd()
	join := filepath.Join(wd, path)
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		join += string(os.PathSeparator)
	}
	return join
}

// getDepth returns the depth of the current path.
func getDepth(givenAbsolutePath, path string) int {
	absoluteSlashCount := strings.Count(givenAbsolutePath, string(os.PathSeparator))
	pathSlashCount := strings.Count(path, string(os.PathSeparator))
	count := pathSlashCount - absoluteSlashCount
	count = count - 1
	if count == -1 {
		count = 0
	}
	return count
}

func (c *Copy) isDestinationDir() bool {
	return strings.HasSuffix(c.GetDestination(), string(os.PathSeparator))
}

func (c *Copy) getCommonPath(path string) string {
	path = getPath(path)
	replaced := strings.Replace(path, c.GetSource(), "", 1)
	// if first character is path separator, remove it
	if strings.HasPrefix(replaced, string(os.PathSeparator)) {
		replaced = strings.TrimPrefix(replaced, string(os.PathSeparator))
	}
	return replaced
}

func (c *CopyOptions) getDirectoryMode() os.FileMode {
	if c.DirectoryMode == nil {
		return constant.DefaultDirPermission
	}
	return *c.DirectoryMode
}

func (c *CopyOptions) getFileMode() os.FileMode {
	if c.FileMode == nil {
		return constant.DefaultFilePermission
	}
	return *c.FileMode
}
