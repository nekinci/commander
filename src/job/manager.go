package job

import (
	"commander/src/specification"
	"fmt"
	"io"
	"os"
)

// Manager is a job manager that keeps track of all jobs and their status.
type Manager struct {
	Name          string
	Version       string
	specification *specification.Specification
	jobs          map[string]*Job
	outputs       map[string]map[string]string // outputs[jobName][commandName] = output
	stdOut        io.Writer
	stdErr        io.Writer
}

func NewManager(specification *specification.Specification) *Manager {
	return NewManagerWithWriters(specification, os.Stdout, os.Stderr)
}

func NewManagerWithWriters(specification *specification.Specification, stdOut io.Writer, stdErr io.Writer) *Manager {
	if specification == nil {
		return nil
	}

	m := &Manager{
		specification: specification,
		jobs:          make(map[string]*Job),
		outputs:       map[string]map[string]string{},
		Name:          specification.Name,
		Version:       specification.Version,
		stdOut:        stdOut,
		stdErr:        stdErr,
	}

	for _, job := range specification.Jobs {
		m.jobs[job.Name] = NewJob(&job, m)
	}

	return m
}

func (m *Manager) Run() {
	for _, job := range m.jobs {
		err := job.Run()
		if err != nil {
			fmt.Printf("Error running job %s: %v\n", job.Name, err)

		}
	}
}
