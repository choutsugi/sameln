package generic

import (
	"LogAgent/universal/error"
)

const (
	TryCloseWithMaxTime = 10
	TrySyncWithMaxTime  = 10
)

type FileUpdateMsg struct {
	FileName    string
	IsUnmarshal bool
	Raw         error.RawErr
}

var ConfigFileUpdateChan = make(chan FileUpdateMsg)
