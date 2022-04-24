package job

import (
	"commander/src/constant"
	"commander/src/job/builtin_jobs"
	"os"
)

type BuiltinJob struct {
	Name   string
	Uses   string
	Params map[string]*interface{}
}

type BuiltinJobRunner interface {
	RunJob() (string, error)
}

func NewBuiltinJob(id string, name, uses string, params map[string]*interface{}) *BuiltinJob {
	return &BuiltinJob{
		Uses:   uses,
		Params: params,
	}
}

func (b *BuiltinJob) Run() (string, error) {
	return NewBuiltinJobRunner(b).RunJob()
}

func getValue[T comparable](v interface{}) *T {
	if v == nil {
		return nil
	}

	return v.(*T)

}

func NewBuiltinJobRunner(job *BuiltinJob) BuiltinJobRunner {
	switch job.Uses {
	case string(constant.COPY_JOB):
		c := builtin_jobs.CopyBuiltinJob{
			Job:       job,
			Source:    (*job.Params["source"]).(string),
			Dest:      (*job.Params["dest"]).(string),
			Recursive: (*job.Params["recursive"]).(bool),
			Depth:     getValue[int](job.Params["depth"]),
			FileMode:  getValue[os.FileMode](job.Params["file_mode"]),
			DirMode:   getValue[os.FileMode](job.Params["dir_mode"]),
		}
		return &c
	}
	return nil
}
