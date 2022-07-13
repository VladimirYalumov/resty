package redis

import (
	"bytes"
	"github.com/go-redis/redis"
)

const EmailAuthCode = "auth_code_"
const EmailAuthCodeCount = "auth_code_count_"

var Client *redis.Client

func InitRedis(host string) {
	Client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})

	_, err := Client.Ping().Result()
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
