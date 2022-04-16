package constant

type InvalidSrcError struct {
	error
	src         string
	description string
}

type InvalidDestError struct {
	error
	dest        string
	description string
}

type EnvVarNotFoundError struct {
	error
	envVarName string
}

func NewEnvVarNotFoundError(envVarName string) *EnvVarNotFoundError {
	return &EnvVarNotFoundError{
		envVarName: envVarName,
	}
}

func NewInvalidSrcError(src, description string) *InvalidSrcError {
	return &InvalidSrcError{src: src, description: description}
}

func NewInvalidDestError(dest, description string) *InvalidDestError {
	return &InvalidDestError{dest: dest, description: description}
}

func (e *InvalidSrcError) Error() string {
	return "Invalid source: " + e.src + " " + e.description
}

func (e *InvalidDestError) Error() string {
	return "Invalid destination: " + e.dest + " " + e.description
}

func (e *EnvVarNotFoundError) Error() string {
	return "Environment variable not found: " + e.envVarName
}
