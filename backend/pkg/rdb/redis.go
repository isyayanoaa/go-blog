package rdb

import (
	"fmt"

	"github.com/isyayanoaa/go-blog/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	c := config.C.Redis
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
	})
}
