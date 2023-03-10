package redis

//Conf 配置信息
type Conf struct {
	Addr
	Database  int  `yaml:"database" toml:"database" json:"database"`       //数据库名
	SSLEnable bool `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"` //是否启用ssl
}

//ClusterConf 集群配置信息
type ClusterConf struct {
	SSLEnable bool   `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"` //是否启用ssl
	AddrList  []Addr `yaml:"addr_list" toml:"addr_list" json:"addr_list"`    //连接地址列表
}

type Addr struct {
	Host     string `yaml:"host" toml:"host" json:"host" default:"localhost"` //主机地址
	Port     int    `yaml:"port" toml:"port" json:"port" default:"6379"`      //端口
	Username string `yaml:"username" toml:"username" json:"username"`         //用户名
	Password string `yaml:"password" toml:"password" json:"password"`         //密码
}
