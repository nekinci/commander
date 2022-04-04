package job

import (
	"commander/src/constant"
	"commander/src/job/builtin_jobs"
)

type BuiltinJob struct {
	Job
	Uses   string
	Params map[string]string
}

func NewBuiltinJob(name, uses string, params map[string]string) *BuiltinJob {
	return &BuiltinJob{
		Job: Job{
			Name: name,
		},
		Uses:   uses,
		Params: params,
	}
}

func (b *BuiltinJob) Run() (string, error) {
	return NewBuiltinJobRunner(b).RunCommand()
}

func NewBuiltinJobRunner(job *BuiltinJob) BuiltinJobRunner {
	switch job.Uses {
	case string(constant.COPY_JOB):
		c := builtin_jobs.CopyBuiltinJob{
			Job:    job,
			Source: job.Params["source"],
			Dest:   job.Params["dest"],
		}
		return &c
	}
	return nil
}
