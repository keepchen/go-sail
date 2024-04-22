package redis

// Conf 配置信息
//
// <yaml example>
//
// redis_conf:
//
// enable = false
//
//	endpoint:
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
// [redis_conf]
//
// enable = false
//
// [redis_conf.endpoint]
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
	Endpoint  `yaml:"endpoint" toml:"endpoint" json:"endpoint"`
	Enable    bool `yaml:"enable" toml:"enable" json:"enable" default:"false"` //是否启用
	Database  int  `yaml:"database" toml:"database" json:"database"`           //数据库名
	SSLEnable bool `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"`     //是否启用ssl
}

// ClusterConf 集群配置信息
//
// <yaml example>
//
// redis_cluster_conf:
//
//		 enable: false
//		 ssl_enable: false
//		 endpoints:
//	   - host: 127.0.0.1
//			port: 6379
//			username: ""
//			password: Mt583611
//	   - host: 127.0.0.1
//			port: 6380
//			username: ""
//			password: Mt583611
//		  - host: 127.0.0.1
//			port: 6381
//			username: ""
//			password: Mt583611
//
// <toml example>
//
// # ::redis cluster配置::
//
// [redis_cluster_conf]
//
// enable = false
//
// ssl_enable = false
//
// [[redis_cluster_conf.endpoints]]
//
// host = "localhost"
//
// username = ""
//
// port = 6379
//
// password = ""
//
// [[redis_cluster_conf.endpoints]]
//
// host = "localhost"
//
// username = ""
//
// port = 6380
//
// password = ""
//
// [[redis_cluster_conf.endpoints]]
//
// host = "localhost"
//
// username = ""
//
// port = 6381
//
// password = ""
type ClusterConf struct {
	Enable    bool       `yaml:"enable" toml:"enable" json:"enable" default:"false"` //是否启用
	SSLEnable bool       `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"`     //是否启用ssl
	Endpoints []Endpoint `yaml:"endpoints" toml:"endpoints" json:"endpoints"`        //连接地址列表
}

type Endpoint struct {
	Host     string `yaml:"host" toml:"host" json:"host" default:"localhost"` //主机地址
	Port     int    `yaml:"port" toml:"port" json:"port" default:"6379"`      //端口
	Username string `yaml:"username" toml:"username" json:"username"`         //用户名
	Password string `yaml:"password" toml:"password" json:"password"`         //密码
}
