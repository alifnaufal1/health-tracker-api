package apperror

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string { return e.Message }