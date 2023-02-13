package http

type Error struct {
	Message    string `json:"message"`
	DevMessage string `json:"dev_message"`
}

func NewErrorMessage(errMessage string, err error) Error {
	return Error{
		Message:    errMessage,
		DevMessage: err.Error(),
	}
}
