package config

type Config struct {
	MysqlUser              string `env:"MYSQL_USER,required"`
	MysqlPassword          string `env:"MYSQL_PASSWORD,required"`
	MysqlDatabase          string `env:"MYSQL_DATABASE,required"`
	MysqlHost              string `env:"MYSQL_HOST,required"`
	Port                   string `env:"PORT, default=8080"`
	NotificationServiceURL string `env:"NOTIFICATION_SERVICE_URL,required"`
	RedisHost              string `env:"REDIS_HOST,required"`
	RedisPassword          string `env:"REDIS_PASSWORD,required"`
	IsMessageProcessing    bool
}

func New() Config {
	return Config{
		IsMessageProcessing: true,
	}
}

func (c *Config) SetMessageProcessing(enable bool) {
	c.IsMessageProcessing = enable
}
