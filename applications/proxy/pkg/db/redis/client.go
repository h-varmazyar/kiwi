package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Configs struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func NewClient(configs *Configs) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:               "",
		Addr:                  fmt.Sprintf("%v:%v", configs.Host, configs.Port),
		ClientName:            "",
		Dialer:                nil,
		OnConnect:             nil,
		Protocol:              0,
		Username:              configs.Username,
		Password:              configs.Password,
		CredentialsProvider:   nil,
		DB:                    configs.DB,
		MaxRetries:            0,
		MinRetryBackoff:       0,
		MaxRetryBackoff:       0,
		DialTimeout:           0,
		ReadTimeout:           0,
		WriteTimeout:          0,
		ContextTimeoutEnabled: false,
		PoolFIFO:              false,
		PoolSize:              0,
		PoolTimeout:           0,
		MinIdleConns:          0,
		MaxIdleConns:          0,
		MaxActiveConns:        0,
		ConnMaxIdleTime:       0,
		ConnMaxLifetime:       0,
		TLSConfig:             nil,
		Limiter:               nil,
		DisableIndentity:      false,
	})
}
