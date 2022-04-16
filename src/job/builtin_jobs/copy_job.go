package builtin_jobs

import (
	"commander/src/constant"
	"commander/src/job"
	"path/filepath"
)

type CopyBuiltinJob struct {
	Job    *job.BuiltinJob
	Source string
	Dest   string
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

func (j *CopyBuiltinJob) RunCommand() (string, error) {

	err := j.validate()
	if err != nil {
		return "", err
	}

	src := j.Source
	dest := j.Dest
	_, _ = src, dest

	sourceAbs, err := filepath.Abs(src)
	_ = sourceAbs

	return "", nil
}
