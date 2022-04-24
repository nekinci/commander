package job

import (
	"commander/src/specification"
	"log"
)

type Job interface {
	Run() (string, error)
}

// validate fatal if the job is invalid
func validate(job *specification.Job, cmd *specification.Command) {

	if job == nil {
		log.Fatal("job is nil")
	}

	if cmd == nil {
		log.Fatal("command is nil")
	}

	if cmd.Uses != nil && cmd.Cmd != nil {
		log.Fatalf("Job %s has a command [%s : %s] that uses both a command and a uses", job.Name, cmd.Id, cmd.Name)
	}
}

func NewJob(specJob *specification.Job, jobCommand *specification.Command) Job {
	validate(specJob, jobCommand)
	if jobCommand.Uses != nil {
		return NewBuiltinJob(jobCommand.Id, jobCommand.Name, *jobCommand.Uses, jobCommand.Params)
	} else if jobCommand.Cmd != nil {
		return NewCommandJob(jobCommand.Id, jobCommand.Name, *jobCommand.Cmd)
	} else {
		log.Fatalf("Job %s has a command [%s : %s] that has no command or uses", specJob.Name, jobCommand.Id, jobCommand.Name)
	}
	return nil
}

type Metadata struct {
	id         string
	name       string
	jobs       []*Job
	definition *specification.Job
}

type Manager struct {
	specification *specification.Specification
	jobs          []*Metadata
}

func registerCommands(metadata *Metadata, job *specification.Job) {
	for _, cmd := range job.Commands {
		r := NewJob(job, &cmd)
		metadata.addJob(&r)
	}
}

func (m *Manager) RegisterJobs(spec *specification.Specification) {
	for _, job := range spec.Jobs {
		m.RegisterJob(&job)
	}
}

func (m *Manager) RegisterJob(job *specification.Job) {
	metadata := &Metadata{
		id:         job.Id,
		name:       job.Name,
		definition: job,
	}
	registerCommands(metadata, job)
	m.jobs = append(m.jobs, metadata)
}

func (m *Metadata) addJob(job *Job) {
	m.jobs = append(m.jobs, job)
}

// NewManager creates a new job manager from a specification. It will register all jobs and commands at the same time.
func NewManager(specification *specification.Specification) *Manager {
	m := &Manager{
		specification: specification,
	}
	m.RegisterJobs(specification)
	return m
}
