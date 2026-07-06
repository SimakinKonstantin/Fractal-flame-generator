package imagegenerator

import "fmt"

type GeneratorError struct {
	msg string
	err error
}

func (e GeneratorError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e GeneratorError) Unwrap() error {
	return e.err
}
