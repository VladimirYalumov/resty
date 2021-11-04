package redis

import (
	"bytes"
	"github.com/go-redis/redis"
)

const REDIS_EMAIL_AUTH_CODE = "auth_code_"
const REDIS_EMAIL_AUTH_CODE_COUNT = "auth_code_count_"

var RedisClient *redis.Client

func InitRedis(host string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func CreateKey(prefix string, postfix string) string {
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(postfix)
	return buffer.String()
}
