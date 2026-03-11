package valkey

// Conf 配置信息
type Conf struct {
	Enable    bool       `yaml:"enable" toml:"enable" json:"enable"`          //是否启用
	Username  string     `yaml:"username" toml:"username" json:"username"`    //用户名
	Password  string     `yaml:"password" toml:"password" json:"password"`    //密码
	Endpoints []Endpoint `yaml:"endpoints" toml:"endpoints" json:"endpoints"` //连接地址列表
}

// Endpoint 节点信息
type Endpoint struct {
	Host string `yaml:"host" toml:"host" json:"host"` //主机地址
	Port int    `yaml:"port" toml:"port" json:"port"` //端口
}
