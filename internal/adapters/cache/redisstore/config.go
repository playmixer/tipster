package redisstore

type Config struct {
	Addr     string `env:"REDIS_ADDRESS"`
	Password string `env:"REDIS_PASSWORD"`
}
