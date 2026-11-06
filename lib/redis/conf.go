package redis

// Conf configuration
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
// # ::Redis configuration::
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
	Enable    bool `yaml:"enable" toml:"enable" json:"enable"`             //Whether to enable
	Database  int  `yaml:"database" toml:"database" json:"database"`       //Database name
	SSLEnable bool `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"` //Whether to enable SSL
}

// ClusterConf cluster configuration
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
// # ::Redis cluster configuration::
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
	Enable    bool       `yaml:"enable" toml:"enable" json:"enable"`             //Whether to enable
	SSLEnable bool       `yaml:"ssl_enable" toml:"ssl_enable" json:"ssl_enable"` //Whether to enable SSL
	Endpoints []Endpoint `yaml:"endpoints" toml:"endpoints" json:"endpoints"`    //Connection address list
}

type Endpoint struct {
	Host     string `yaml:"host" toml:"host" json:"host"`             //Host address
	Port     int    `yaml:"port" toml:"port" json:"port"`             //Port
	Username string `yaml:"username" toml:"username" json:"username"` //Username
	Password string `yaml:"password" toml:"password" json:"password"` //Password
}
