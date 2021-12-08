package blunder

func New(code uint64, err error) *Errors {
	return &Errors{
		Code: code,
		Msg:  blunderMsg[code],
		Err:  err,
	}
}

func NewWithSuccess() *Errors {
	return &Errors{
		Code: CODE_SUCCESS,
		Msg:  blunderMsg[CODE_SUCCESS],
		Err:  nil,
	}
}

func GetMsg(code uint64) string {
	return blunderMsg[code]
}
