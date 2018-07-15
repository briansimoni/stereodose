package controllers

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

// StatusError is handy for when you want to return something other than 500 internal server error
type statusError struct {
	error
	Message string
	Code    int
}

func (e *statusError) Error() string {
	return e.Message
}

// Status returns the http response code
func (e *statusError) Status() int {
	return e.Code
}
