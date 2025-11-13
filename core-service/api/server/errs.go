package server

import "fmt"

type AddError struct {
	Number string
	Err    error
}

func (a AddError) Error() string {
	return fmt.Sprintf("%s: %s", a.Number, a.Err.Error())
}
