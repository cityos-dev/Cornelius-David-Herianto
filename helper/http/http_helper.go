package http

// Error represents error message that will be reported to the client
type Error struct {
	Message    string `json:"message"`
	DevMessage string `json:"dev_message"`
}

// NewErrorMessage returned new Error object
func NewErrorMessage(errMessage string, err error) Error {
	return Error{
		Message:    errMessage,
		DevMessage: err.Error(),
	}
}
