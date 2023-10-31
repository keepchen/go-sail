package nats

// Conf 配置
//
// <yaml example>
//
// nats:
//
//	servers:
//	  - "nats://127.0.0.1:4222"
//	username: admin
//	password: changeMe
//
// <toml example>
//
// # ::nats配置::
//
// [nats]
//
// # 服务实例地址
//
// servers = ["nats://127.0.0.1:4222"]
//
// # 认证用户名
//
// username = "admin"
//
// # 认证密码
//
// password = "changeMe"
type Conf struct {
	Servers  []string `yaml:"servers" toml:"servers" json:"servers"`    //服务实例列表
	Username string   `yaml:"username" toml:"username" json:"username"` //用户名
	Password string   `yaml:"password" toml:"password" json:"password"` //密码
}
