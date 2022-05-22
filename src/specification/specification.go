package specification

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

/*
```yaml
	version: v1
	name: "example"
	on:
		change:
			files: []
		save:
			files: []
		delete:
			files: []
		create:
			directory: []
	jobs:
		job_one:
			name: "job_one"
			commands:
				- id: "step-1"
				  uses: "copy"
				  params:
					source: "source"
					destination: "destination"
					recursive: true
					depth: 1
				- id: "step-2"
				  cmd: "echo hello world"
		job_two:
			name: "job_two"
			commands:
				- id: "step-1"
				  uses: "copy"
				  params:
					source: "source"
					destination: "destination"
					recursive: true
					depth: 1
				- id: "step-2"
				  cmd: "echo hello world" // windows or linux command
```


*/

type Change struct {
	Files []string `yaml:"files"`
}

type Create struct {
	Directories []string `yaml:"directory"`
}

type Timer struct {
	Cron         string `yaml:"cron"`
	Milliseconds int    `yaml:"milliseconds"`
	Minutes      int    `yaml:"minutes"`
	MaxRetries   int    `yaml:"max_retries"` // -1 for infinite
}

type On struct {
	Change Change `yaml:"change"`
	Save   Change `yaml:"save"`
	Delete Change `yaml:"delete"`
	Create Create `yaml:"create"`
	Timer  Timer  `yaml:"timer"`
}

type Command struct {
	Params map[string]*interface{} `yaml:"params"`
	Uses   *string                 `yaml:"uses"`
	Cmd    *string                 `yaml:"cmd"`
	Name   string                  `yaml:"name"`
	Id     string                  `yaml:"id"`
}

type Job struct {
	Id       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

type Specification struct {
	Version string         `yaml:"version"`
	Name    string         `yaml:"name"`
	On      On             `yaml:"on"`
	Jobs    map[string]Job `yaml:"jobs"`
}

func LoadYaml(path string) (Specification, error) {
	var spec Specification
	file, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return spec, err2
	}

	err := yaml.Unmarshal(file, &spec)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

// LoadFromFile loads a specification from a file
func LoadFromFile(file string) (*Specification, error) {
	spec, err := LoadYaml(file) // if we want to change file format, we can do it here
	if err != nil {
		return nil, err
	}
	return &spec, nil
}
