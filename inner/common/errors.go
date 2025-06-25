package common

type RequestValidationError struct {
	Message string
}

func (err RequestValidationError) Error() string {
	return err.Message
}

type AlreadyExistsError struct {
	Message string
}

func (err AlreadyExistsError) Error() string {
	return err.Message
}

type NotFoundError struct {
	Message string
}

func (err NotFoundError) Error() string {
	return err.Message
}
