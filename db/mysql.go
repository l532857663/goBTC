package db

type Mysql struct {
	Config       string `mapstructure:"config"         json:"config"        yaml:"config"`         // 配置信息
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns"  yaml:"max-idle-conns"` // 最大空闲连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns"  yaml:"max-open-conns"` // 最大打开连接数
	LogMode      bool   `mapstructure:"log-mode"       json:"logMode"       yaml:"log-mode"`       // 日志开关
	LogZap       string `mapstructure:"log-zap"        json:"logZap"        yaml:"log-zap"`        // 日志级别
	Prefix       string `mapstructure:"prefix"         json:"prefix"        yaml:"prefix"`         // 前缀
	Path         string `mapstructure:"path"           json:"path"          yaml:"path"`           // 连接地址
	Dbname       string `mapstructure:"db-name"        json:"dbname"        yaml:"db-name"`        // 数据库名
	Username     string `mapstructure:"username"       json:"username"      yaml:"username"`       // 账户名
	Password     string `mapstructure:"password"       json:"password"      yaml:"password"`       // 密码
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}
