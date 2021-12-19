package bundles

// FileUpdateMsg 文件更新channel通信包
type FileUpdateMsg struct {
	FileName    string
	IsUnmarshal bool
}

// ConfigFileUpdateChan 配置文件更新通知channel
var ConfigFileUpdateChan = make(chan FileUpdateMsg)
