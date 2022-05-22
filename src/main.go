package main

import "C"
import (
	"commander/src/job"
	"commander/src/listener"
	"commander/src/specification"
)

func main() {

	ss := "touch a.txt"
	ss1 := "echo 'hello world!'"
	ss2 := "echo naber hello"
	ss3 := "git status"
	s := specification.Specification{
		Version: "v1",
		Name:    "Niyazi Test",

		On: specification.On{
			Change: specification.Change{
				Files: []string{"./src"},
			},
			Save: specification.Change{
				Files: []string{"./src"},
			},
			Delete: specification.Change{},
			Create: specification.Create{},
			Timer:  specification.Timer{},
		},
		Jobs: map[string]specification.Job{
			"job": {
				Id:   "new-job-1",
				Name: "new-job-1",
				Commands: []specification.Command{
					{Params: nil, Uses: nil, Cmd: &ss, Name: "c1", Id: "c1"},
					{Params: nil, Uses: nil, Cmd: &ss1, Name: "c2", Id: "c2"},
					{Params: nil, Uses: nil, Cmd: &ss2, Name: "c2", Id: "c2"},
					{Params: nil, Uses: nil, Cmd: &ss3, Name: "c2", Id: "c2"},
				},
			},
		},
	}

	manager := job.NewManager(&s)
	l := listener.New(&s, manager)

	l.Listen()
}
