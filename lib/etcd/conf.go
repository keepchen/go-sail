package etcd

import "crypto/tls"

type Conf struct {
	Enable    bool        `yaml:"enable" toml:"enable" json:"enable"`          //Whether to enable
	Endpoints []string    `yaml:"endpoints" toml:"endpoints" json:"endpoints"` //Address list, e.g.: localhost:2379
	Username  string      `yaml:"username" toml:"username" json:"username"`    //Account
	Password  string      `yaml:"password" toml:"password" json:"password"`    //Password
	Timeout   int         `yaml:"timeout" toml:"timeout" json:"timeout"`       //Connection timeout (milliseconds), default: 10000ms
	Tls       *tls.Config `yaml:"-" toml:"-" json:"-"`                         //TLS configuration
}
