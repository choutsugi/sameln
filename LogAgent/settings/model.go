package settings

type ConfigType struct {
	App      *AppConfigType      `mapstructure:"app"`
	Log      *LogConfigType      `mapstructure:"log"`
	Etcd     *EtcdConfigType     `mapstructure:"etcd"`
	Kafaka   *KafkaConfigType    `mapstructure:"kafka"`
	Nsq      *NsqConfigType      `mapstructure:"nsq"`
	InfluxDB *InfluxDbConfigType `mapstructure:"influxDB"`
}

type AppConfigType struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
}

type LogConfigType struct {
	Level      string `mapstructure:"level"`
	Type       string `mapstructure:"type"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type EtcdConfigType struct {
	Addr       string `mapstructure:"addr"`
	Port       int    `mapstructure:"port"`
	CollectKey string `mapstructure:"collect_key"`
}

type KafkaConfigType struct {
	Addr     string `mapstructure:"addr"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ChanSize int    `mapstructure:"chan_size"`
}

type NsqConfigType struct {
	Addr     string `mapstructure:"addr"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ChanSize int    `mapstructure:"chan_size"`
}

type InfluxDbConfigType struct {
	Addr      string `mapstructure:"addr"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Precision string `mapstructure:"precision"`
}

const (
	ModeRelease = "release"
	ModeDevelop = "develop"
)
