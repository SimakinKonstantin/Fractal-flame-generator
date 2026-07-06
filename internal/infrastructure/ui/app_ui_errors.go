package ui

import "fmt"

type InputError struct {
	msg string
	err error
}

func (e InputError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e InputError) Unwrap() error {
	return e.err
}

type OutputError struct {
	msg string
	err error
}

func (e OutputError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e OutputError) Unwrap() error {
	return e.err
}

type FileError struct {
	msg string
	err error
}

func (e FileError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e FileError) Unwrap() error {
	return e.err
}
