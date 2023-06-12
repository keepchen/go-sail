package redis

// Conf 配置信息
//
// <yaml example>
//
// redis:
//
//	addr:
//	  host: "127.0.0.1"
//	  port: 6379
//	  username: ""
//	  password: ""
//	database: 0
//	ssl_enable: false
//
// <toml example>
//
// # ::redis配置::
//
// [redis]
//
// host = "localhost"
//
// username = ""
//
// port = 6379
//
// password = ""
//
// database = 0
//
// ssl_enable = false
type Conf struct {
	Addr
	Database  int  `yaml:"database" toml:"database" json:"database"`       //数据库名
	SSLEnable bool `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"` //是否启用ssl
}

// ClusterConf 集群配置信息
//
// <yaml example>
//
// redis_cluster:
//
//	ssl_enable: false
//	addr_list:
//	  - host: 127.0.0.1
//	    port: 6379
//	    username: ""
//	    password: Mt583611
//	  - host: 127.0.0.1
//	    port: 6380
//	    username: ""
//	    password: Mt583611
//	  - host: 127.0.0.1
//	    port: 6381
//	    username: ""
//	    password: Mt583611
//
// <toml example>
//
// # ::redis cluster配置::
//
// [redis_cluster]
//
// ssl_enable = false
//
// [[redis_cluster.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6379
//
// password = ""
//
// [[redis_cluster.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6380
//
// password = ""
//
// [[redis_cluster.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6381
//
// password = ""
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
