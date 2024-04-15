package email

// Conf 配置
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
	Workers               int    `yaml:"workers" toml:"workers" json:"workers"`                                                 //协程数量
	WorkerThrottleSeconds int    `yaml:"worker_throttle_seconds" toml:"worker_throttle_seconds" json:"worker_throttle_seconds"` //每个协程内发送间隔，单位秒
	Host                  string `yaml:"host" toml:"host" json:"host"`                                                          //邮件服务器域名
	Port                  int    `yaml:"port" toml:"port" json:"port"`                                                          //邮件服务器端口
	Username              string `yaml:"username" toml:"username" json:"username"`                                              //邮件服务登录账号
	Password              string `yaml:"password" toml:"password" json:"password"`                                              //邮件服务登录密码
	From                  string `yaml:"from" toml:"from" json:"from"`                                                          //发送人邮箱
	Subject               string `yaml:"subject" toml:"subject" json:"subject"`                                                 //邮件主题
	Params                struct {
		Variables []string `yaml:"variables" toml:"variables" json:"variables"` //替换变量
	} `yaml:"params" toml:"params" json:"params"` //其他参数
}
