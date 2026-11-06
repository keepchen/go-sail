package nats

// Conf configuration
//
// <yaml example>
//
// nats_conf:
//
//	endpoints:
//	  - "nats://127.0.0.1:4222"
//	username: admin
//	password: changeMe
//
// <toml example>
//
// # ::NATS configuration::
//
// [nats_conf]
//
// # Service instance address
//
// endpoints = ["nats://127.0.0.1:4222"]
//
// # Authentication username
//
// username = "admin"
//
// # Authentication password
//
// password = "changeMe"
type Conf struct {
	Enable    bool     `yaml:"enable" toml:"enable" json:"enable"`          //Whether to enable
	Endpoints []string `yaml:"endpoints" toml:"endpoints" json:"endpoints"` //Service instance list
	Username  string   `yaml:"username" toml:"username" json:"username"`    //Username
	Password  string   `yaml:"password" toml:"password" json:"password"`    //Password
}
