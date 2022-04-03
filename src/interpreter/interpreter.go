package interpreter

type Interpreter interface {
	Interpret(string) (string, error)
}
