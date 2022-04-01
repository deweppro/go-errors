package errors

import "fmt"

type errMessage struct {
	cause   error
	message string
	trace   string
}

func New(message string) error {
	return &errMessage{message: message}
}

func (v *errMessage) Error() string {
	switch true {
	case len(v.message) > 0 && v.cause != nil:
		return v.message + ": " + v.cause.Error() + v.trace
	case v.cause != nil:
		return v.cause.Error() + v.trace
	}
	return v.message + v.trace
}

func (v *errMessage) Cause() error {
	return v.cause
}

func (v *errMessage) Unwrap() error {
	return v.cause
}

func (v *errMessage) WithTrace() {
	v.trace = tracing()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func WrapMessage(cause error, message string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	var err0 *errMessage
	if len(args) == 0 {
		err0 = &errMessage{
			cause:   cause,
			message: message,
		}
	} else {
		err0 = &errMessage{
			cause:   cause,
			message: fmt.Sprintf(message, args...),
		}
	}
	err0.WithTrace()
	return err0
}

func Wrap(msg ...error) error {
	if len(msg) == 0 {
		return nil
	}
	var err0 *errMessage
	for _, v := range msg {
		if v == nil {
			continue
		}
		if err0 == nil {
			err0 = &errMessage{cause: v}
			continue
		}
		err0 = &errMessage{
			cause:   v,
			message: err0.Error(),
		}
	}
	return err0
}
