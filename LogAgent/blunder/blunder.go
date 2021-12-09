package blunder

var (
	success = &Error{
		State: StatusSuccess,
		Code:  CodeSysSuccess,
		Msg:   blunderMsg[CodeSysSuccess],
		Err:   nil,
	}
)

func NewError(code uint64, err error) *Error {
	return &Error{
		State: StatusFailed,
		Code:  code,
		Msg:   blunderMsg[code],
		Err:   err,
	}
}

func Success() *Error {
	return success
}

func NewSuccess(code uint64) *Error {
	return &Error{
		State: StatusSuccess,
		Code:  code,
		Msg:   blunderMsg[code],
		Err:   nil,
	}
}

func GetMsg(code uint64) string {
	return blunderMsg[code]
}
