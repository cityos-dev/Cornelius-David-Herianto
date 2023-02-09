package http

type HTTPError struct {
	Message string `json:"message"`
}

func NewHTTPErrorMessage(err error) HTTPError {
	return HTTPError{
		Message: err.Error(),
	}
}
