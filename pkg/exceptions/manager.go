package exceptions

type ErrManager struct {
	Message string
}

func (e ErrManager) Error() string {
	return e.Message
}
