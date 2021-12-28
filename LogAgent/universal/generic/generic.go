package generic

const (
	TryCloseWithMaxTime = 10 // 最大尝试关闭次数
	TrySyncWithMaxTime  = 10
)

// FileUpdateMsg 文件更新channel通信包
type FileUpdateMsg struct {
	FileName    string
	IsUnmarshal bool
}

// ConfigFileUpdateChan 配置文件更新通知channel
var ConfigFileUpdateChan = make(chan FileUpdateMsg)
