package qbapi

import "fmt"

type QError struct {
	code int
	err  error
}

func (q *QError) Error() string {
	return fmt.Sprintf("QError:[code:%d, err:%+v]", q.code, q.err)
}

func (q *QError) Code() int {
	return q.code
}

func (q *QError) Err() error {
	return q.err
}

func (q *QError) RootCause() (int, error) {
	err := q.err
	code := q.code
	for {
		if err == nil {
			break
		}
		e, ok := err.(*QError)
		if !ok {
			break
		}
		err = e.err
		code = e.code
	}
	return code, err
}

func NewError(code int, err error) *QError {
	return &QError{code: code, err: err}
}

func NewMsgError(code int, msg string) *QError {
	return NewError(code, fmt.Errorf("%s", msg))
}

func RootCause(err error) (int, error) {
	e, ok := err.(*QError)
	if !ok {
		return ErrUnknown, err
	}
	return e.RootCause()
}

type StatusCodeErr struct {
	code int
}

func NewStatusCodeErr(code int) *StatusCodeErr {
	return &StatusCodeErr{code: code}
}

func (s *StatusCodeErr) Code() int {
	return s.code
}

func (s *StatusCodeErr) Error() string {
	return fmt.Sprintf("StatusCodeErr:[code:%d]", s.code)
}
