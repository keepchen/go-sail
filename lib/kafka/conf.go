package kafka

import "crypto/tls"

type Conf struct {
	Enable       bool        `yaml:"enable" toml:"enable" json:"enable"`                   //Whether to enable
	Endpoints    []string    `yaml:"endpoints" toml:"endpoints" json:"endpoints"`          //Address list, e.g.: localhost:9092
	SASLAuthType string      `yaml:"SASLAuthType" toml:"SASLAuthType" json:"SASLAuthType"` //Authentication encryption method: plain, sha256, sha512
	Username     string      `yaml:"username" toml:"username" json:"username"`             //Account
	Password     string      `yaml:"password" toml:"password" json:"password"`             //Password
	Timeout      int         `yaml:"timeout" toml:"timeout" json:"timeout"`                //Connection timeout (milliseconds), default: 10000ms
	Tls          *tls.Config `yaml:"-" toml:"-" json:"-"`                                  //TLS configuration
}
