package job

import (
	"commander/src/specification"
	"errors"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

type collectableWriter struct {
	writer io.Writer
	buffer []byte
}

func (cw *collectableWriter) Write(p []byte) (n int, err error) {
	cw.buffer = append(cw.buffer, p...)
	return cw.writer.Write(p)
}

type Job struct {
	Id         string
	Name       string
	definition *specification.Job
	Commands   []*Command
	manager    *Manager
}

type Command struct {
	Id         string
	Name       string
	definition *specification.Command
	Cmd        *string
	job        *Job
	// Uses and Params not implemented yet may be added later if needed
}

func (c *Command) Run() (string, error) {

	if c.Cmd == nil {
		return "", errors.New("command not defined")
	}
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		command = exec.Command("powershell.exe")
	} else {
		command = exec.Command("/bin/bash")
	}

	command.Stdin = strings.NewReader(*c.Cmd)

	out := collectableWriter{writer: c.job.manager.stdOut}
	errOut := collectableWriter{writer: c.job.manager.stdErr}
	command.Stdout = &out
	command.Stderr = &errOut

	err := command.Run()
	if err != nil {
		return string(errOut.buffer), err
	}
	return string(out.buffer), nil
}

func NewJob(job *specification.Job, m *Manager) *Job {
	j := &Job{
		Id:         job.Id,
		Name:       job.Name,
		definition: job,
		manager:    m,
	}

	j.Commands = make([]*Command, len(job.Commands))
	for i, c := range job.Commands {
		j.Commands[i] = &Command{
			Id:         c.Id,
			Name:       c.Name,
			Cmd:        c.Cmd,
			definition: &c,
			job:        j,
		}
	}

	return j
}

func (j *Job) Run() error {

	if j.manager.outputs[j.Name] == nil {
		j.manager.outputs[j.Name] = make(map[string]string)
	}

	for _, c := range j.Commands {
		out, err := c.Run()
		if err != nil {
			if out != "" {
				j.manager.outputs[j.Name][c.Name] = out
			}
			return err
		}

		j.manager.outputs[j.Name][c.Name] = out
	}

	return nil
}
