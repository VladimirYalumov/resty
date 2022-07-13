package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"resty/mail"
	"resty/orm"
	"resty/redis"
	"resty/routing"
)

type Conf struct {
	DbHost        string `yaml:"db_host"`
	DbName        string `yaml:"db_name"`
	DbUser        string `yaml:"db_user"`
	DbPassword    string `yaml:"db_password"`
	DbPort        string `yaml:"db_port"`
	SmtpServer    string `yaml:"smtp_server"`
	EmailUser     string `yaml:"email_user"`
	EmailPassword string `yaml:"email_password"`
	RedisHost     string `yaml:"redis_host"`
}

var config Conf

func (c *Conf) GetConf() *Conf {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}
	return c
}

func Init() {
	config.GetConf()
	orm.InitDB(config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	mail.InitMailConfig(config.SmtpServer, config.EmailUser, config.EmailPassword)
	redis.InitRedis(config.RedisHost)
	routing.InitErrors()
}

func Close() {
	orm.CloseDB()
}
