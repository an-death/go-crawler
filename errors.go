package main

import "fmt"

type DoesNotMatchError struct {
	expected interface{}
	actual interface{}
}

func (err *DoesNotMatchError) Error() string {
	return fmt.Sprintf(
		"response codes doesn't match. Actual: %v  Expected: %v",
		err.actual, err.expected)
}

