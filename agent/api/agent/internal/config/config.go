package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	CAFile string
	Debug  struct {
		Enable bool
	}
}
