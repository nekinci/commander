package job

type IJob interface {
	Run() (string, error)
}

type BuiltinJobRunner interface {
	RunCommand() (string, error)
}

type Job struct {
	Name string
}

type FileJob struct {
	Job
	Directory string
}
