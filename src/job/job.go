package job

type IJob interface {
	Run() error
}

type Job struct {
	Name string
}

type BuiltinJob struct {
	Job
	Params map[string]string
}

type FileJob struct {
	Job
	Directory string
}
