package hook

type AddError struct {
	Status  string `json:"status"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e AddError) Error() string {
	return e.Message
}
