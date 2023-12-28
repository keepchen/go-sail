package kafka

import "crypto/tls"

type Conf struct {
	AddrList     []string    `yaml:"addrList" toml:"addrList" json:"addrList"`             //地址列表,如: localhost:9092
	SASLAuthType string      `yaml:"SASLAuthType" toml:"SASLAuthType" json:"SASLAuthType"` //认证加密方式：plain、sha256、sha512
	Username     string      `yaml:"username" toml:"username" json:"username"`             //账号
	Password     string      `yaml:"password" toml:"password" json:"password"`             //密码
	Timeout      int         `yaml:"timeout" toml:"timeout" json:"timeout"`                //连接超时时间（毫秒）默认10000ms
	Tls          *tls.Config `yaml:"-" toml:"-" json:"-"`                                  //tls配置//tls配置
}
