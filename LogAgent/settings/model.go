package settings

type ConfigType struct {
	App    *app    `mapstructure:"app"`
	Log    *log    `mapstructure:"log"`
	Etcd   *etcd   `mapstructure:"etcd"`
	Kafaka *kafka  `mapstructure:"kafka"`
	Nsq    *nsq    `mapstructure:"nsq"`
	Influx *influx `mapstructure:"influx"`
}

type app struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
}

type log struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type etcd struct {
	Addr       string `mapstructure:"addr"`
	Port       int    `mapstructure:"port"`
	CollectKey string `mapstructure:"collect_key"`
}

type kafka struct {
	Addr     string `mapstructure:"addr"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ChanSize int    `mapstructure:"chan_size"`
}

type nsq struct {
	Addr     string `mapstructure:"addr"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ChanSize int    `mapstructure:"chan_size"`
}

type influx struct {
	Addr      string `mapstructure:"addr"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Precision string `mapstructure:"precision"`
}
