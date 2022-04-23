package validator

type Validator interface {
	Validate() error
	ShouldRunAtThisPoint() bool
}
