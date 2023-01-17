package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		Host     string
		Port     int
		Account  string
		Password string
		Charset  string
		Name     string
	}
	Redis struct {
		Host string
		Port int
	}
}
