package server

import "strings"

type AddError struct {
	Errs []string
}

func (a *AddError) Error() string {
	return strings.Join(a.Errs, ";\n")
}
