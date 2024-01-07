package exceptions

type ErrRemoveNameSuffix struct {
}

func (e ErrRemoveNameSuffix) ParseFailed() bool {
	return true
}

func (e ErrRemoveNameSuffix) Error() string {
	return "处理名称失败"
}
