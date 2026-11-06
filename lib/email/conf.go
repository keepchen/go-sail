package email

// Conf configuration
//
// <yaml example>
//
// email_conf:
//
//	workers:
//	work_throttle_seconds: 5
//	host:
//	port:
//	username:
//	password:
//	from:
//	subject:
//	params:
//		variables:
//
// <toml example>
//
// [email_conf]
//
// workers =
//
// work_throttle_seconds = 5
//
// host =
//
// port =
//
// username =
//
// password =
//
// from =
//
// subject =
//
// [email_conf.params]
//
// variables =
type Conf struct {
	Workers               int    `yaml:"workers" toml:"workers" json:"workers"`                                                 //Number of goroutines
	WorkerThrottleSeconds int    `yaml:"worker_throttle_seconds" toml:"worker_throttle_seconds" json:"worker_throttle_seconds"` //Send interval within each goroutine, in seconds
	Host                  string `yaml:"host" toml:"host" json:"host"`                                                          //Mail server domain
	Port                  int    `yaml:"port" toml:"port" json:"port"`                                                          //Mail server port
	Username              string `yaml:"username" toml:"username" json:"username"`                                              //Mail service login account
	Password              string `yaml:"password" toml:"password" json:"password"`                                              //Mail service login password
	From                  string `yaml:"from" toml:"from" json:"from"`                                                          //Sender email address
	Subject               string `yaml:"subject" toml:"subject" json:"subject"`                                                 //Email subject
	Params                struct {
		Variables []string `yaml:"variables" toml:"variables" json:"variables"` //Replacement variables
	} `yaml:"params" toml:"params" json:"params"` //Other parameters
}
