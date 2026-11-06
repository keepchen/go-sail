package valkey

// Conf configuration
type Conf struct {
	Enable    bool       `yaml:"enable" toml:"enable" json:"enable"`          //Whether to enable
	Username  string     `yaml:"username" toml:"username" json:"username"`    //Username
	Password  string     `yaml:"password" toml:"password" json:"password"`    //Password
	Endpoints []Endpoint `yaml:"endpoints" toml:"endpoints" json:"endpoints"` //Connection address list
}

// Endpoint node information
type Endpoint struct {
	Host string `yaml:"host" toml:"host" json:"host"` //Host address
	Port int    `yaml:"port" toml:"port" json:"port"` //Port
}
