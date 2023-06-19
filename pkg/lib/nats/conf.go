package nats

// Conf 配置
type Conf struct {
	Servers  []string `yaml:"servers" toml:"servers" json:"servers"`    //服务实例列表
	Username string   `yaml:"username" toml:"username" json:"username"` //用户名
	Password string   `yaml:"password" toml:"password" json:"password"` //密码
}
