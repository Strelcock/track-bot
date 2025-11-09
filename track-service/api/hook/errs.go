package hook

type HookAddError struct {
	Status  string `json:"status"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e HookAddError) Error() string {
	return e.Message
}
