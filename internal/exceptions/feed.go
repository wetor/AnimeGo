package exceptions

type ErrFeed struct {
	Message string
}

func (e ErrFeed) Error() string {
	return e.Message
}
