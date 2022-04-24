package job

import (
	"os"
	"os/exec"
	"strings"
)

type CommandJob struct {
	Command     string
	CommandArgs []string
}

func NewCommandJob(id, name, command string) *CommandJob {
	return &CommandJob{
		Command:     command,
		CommandArgs: determineCommand(command),
	}
}

func determineCommand(c string) []string {
	args := strings.Split(c, " ")
	return args
}

func (c *CommandJob) Cmd() string {
	return c.CommandArgs[0]
}

func (c *CommandJob) Args() []string {
	return c.CommandArgs[1:]
}

func (c *CommandJob) Run() (string, error) {
	cmd := exec.Command(c.Cmd(), c.Args()...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(out), err
}
