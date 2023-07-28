package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DBUri          string
	MaxConnections int
	EXPIRETIME     int
	Etcd           struct {
		Hosts []string
		Key   struct {
			Rest     string
			Ws       string
			SocketIO string
		}
	}
}
