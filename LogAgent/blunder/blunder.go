package blunder

func New(code uint64, err *error) *Errors {
	return &Errors{
		Msg: blunderMsg[code],
		Err: err,
	}
}
