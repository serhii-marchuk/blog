package configs

import "fmt"

type RedisCfg struct {
	HOST     string `env:"REDIS_HOST"`
	PORT     int    `env:"REDIS_PORT"`
	PASSWORD string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

func (r *RedisCfg) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.HOST, r.PORT)
}
