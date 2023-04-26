package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

// Client 客户端
type Client struct {
	cli  *smtp.Client
	auth smtp.Auth
	conf Conf
}

// New 新建smtp实例
func New(conf Conf) *Client {
	return &Client{
		cli:  initSmtpClient(conf),
		auth: authorize(conf),
		conf: conf,
	}
}

// SendMailWithTLS 使用tls方式发送
func (c *Client) SendMailWithTLS(to []string, msg []byte) error {
	if c.cli == nil {
		return errors.New("smtp client is nil")
	}
	if ok, _ := c.cli.Extension("AUTH"); ok {
		if err := c.cli.Auth(c.auth); err != nil {
			return err
		}
	}
	if err := c.cli.Mail(c.conf.Username); err != nil {
		return err
	}
	for _, addr := range to {
		if err := c.cli.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.cli.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return c.cli.Quit()
}

// 初始化smtp客户端
func initSmtpClient(conf Conf) *smtp.Client {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("initialize smtp client failed:", err.Error())
		return nil
	}

	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil
		//panic(err)
	}

	return client
}

// 授权
func authorize(conf Conf) smtp.Auth {
	return smtp.PlainAuth("", conf.Username, conf.Password, conf.Host)
}
