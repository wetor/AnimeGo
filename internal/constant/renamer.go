package constant

const (
	RenameStateError = iota - 1
	RenameStateStart
	RenameStateSeeding
	RenameStateComplete
	RenameStateEnd
)

const (
	AllRenameStateError = iota - 1
	AllRenameStateStart
	AllRenameStateIncomplete
	AllRenameStateComplete
)

const (
	RenameStateChanCap = 5
	RenameMaxErrCount  = 3
)
