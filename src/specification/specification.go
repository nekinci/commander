package specification

type Change struct {
	Files []string `yaml:"files"`
}

type Create struct {
	Directories []string `yaml:"directory"`
}

type On struct {
	Change Change `yaml:"change"`
	Save   Change `yaml:"save"`
	Delete Change `yaml:"delete"`
	Create Create `yaml:"create"`
}

type Command struct {
	Uses string `yaml:"uses"`
	Cmd  string `yaml:"cmd"`
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type Job struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

type Specification struct {
	Version string         `yaml:"version"`
	Name    string         `yaml:"name"`
	On      On             `yaml:"on"`
	Jobs    map[string]Job `yaml:"jobs"`
}
