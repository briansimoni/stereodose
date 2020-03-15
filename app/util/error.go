package util

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

// StatusError is handy for when you want to return something other than 500 internal server error
type StatusError struct {
	error
	Message string
	Code    int
}

func (e *StatusError) Error() string {
	return e.Message
}

// Status returns the http response code
func (e *StatusError) Status() int {
	return e.Code
}
