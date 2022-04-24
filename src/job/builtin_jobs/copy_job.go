package builtin_jobs

import (
	"commander/src/constant"
	"commander/src/job"
	"commander/src/osutil"
	"os"
)

type CopyBuiltinJob struct {
	Job       *job.BuiltinJob
	Source    string
	Dest      string
	Recursive bool
	Depth     *int
	FileMode  *os.FileMode
	DirMode   *os.FileMode
}

func (j *CopyBuiltinJob) validate() error {
	src := j.Source
	dest := j.Dest

	if src == "" {
		return constant.NewInvalidSrcError(src, "Source cannot be empty!")
	}

	if dest == "" {
		return constant.NewInvalidDestError(dest, "Destination cannot be empty!")
	}

	if src == dest {
		return constant.NewInvalidDestError(dest, "Source and destination cannot be the same!")
	}

	if src == "/" {
		return constant.NewInvalidSrcError(src, "Source cannot be root directory!")
	}

	if dest == "/" {
		return constant.NewInvalidDestError(dest, "Destination cannot be root directory!")
	}

	return nil

}

func (j *CopyBuiltinJob) RunJob() (string, error) {

	err := j.validate()
	if err != nil {
		return "", err
	}

	src := j.Source
	dest := j.Dest
	_, _ = src, dest

	options := osutil.CopyOptions{
		Recursive:     j.Recursive,
		Depth:         j.Depth,
		FileMode:      j.FileMode,
		DirectoryMode: j.DirMode,
		FilterFunc:    nil, // TODO: implement filter func by incoming filter
	}
	c := osutil.NewCopy(src, dest, &options)
	return "", c.Copy()
}
