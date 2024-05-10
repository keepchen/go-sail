package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/segmentio/kafka-go/sasl/scram"

	kafkaLib "github.com/segmentio/kafka-go"
)

// Instance 连接实例
type Instance struct {
	Connections []*kafkaLib.Conn
	Writer      *kafkaLib.Writer
	Reader      *kafkaLib.Reader
}

var instance *Instance

// GetConnections 获取连接
//
// 调用此方法前需要先调用 Init 进行初始化
func GetConnections() []*kafkaLib.Conn {
	return instance.Connections
}

// GetWriter 获取写实例
//
// 调用此方法前需要先调用 Init 进行初始化
func GetWriter() *kafkaLib.Writer {
	return instance.Writer
}

// GetReader 获取读实例
//
// 调用此方法前需要先调用 Init 进行初始化
func GetReader() *kafkaLib.Reader {
	return instance.Reader
}

// GetInstance 获取完整实例
//
// 包含连接、读实例、写实例
func GetInstance() *Instance {
	return instance
}

// Init 初始化连接
//
// 该方法会初始化连接、读实例、写实例
//
// 初始化后，可调用 GetInstance 方法获取完整实例
func Init(conf Conf, topic, groupID string) {
	InitConnections(conf)
	InitWriter(conf, topic)
	InitReader(conf, topic, groupID)
}

// New 初始化连接
//
// 该方法会初始化连接、读实例、写实例
func New(conf Conf, topic, groupID string) (connections []*kafkaLib.Conn,
	writer *kafkaLib.Writer, wErr error,
	reader *kafkaLib.Reader, rErr error) {
	connections = NewConnections(conf)
	writer, wErr = NewWriter(conf, topic)
	reader, rErr = NewReader(conf, topic, groupID)

	return
}

// InitConnections 初始化连接
func InitConnections(conf Conf) {
	var connections = make([]*kafkaLib.Conn, 0, len(conf.Endpoints))
	for _, addr := range conf.Endpoints {
		var (
			conn *kafkaLib.Conn
			err  error
		)
		//无账号密码连接
		if len(conf.Username) == 0 && len(conf.Password) == 0 {
			conn, err = kafkaLib.Dial("tcp", addr)
		} else {
			mechanism, mErr := getMechanism(conf)
			if mErr != nil {
				panic(mErr)
			}
			if conf.Timeout < 1 {
				conf.Timeout = 10000
			}
			dialer := &kafkaLib.Dialer{
				Timeout:       time.Duration(conf.Timeout) * time.Millisecond,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
			conn, err = dialer.DialContext(context.Background(), "tcp", addr)
		}
		if err != nil {
			panic(err)
		}
		connections = append(connections, conn)
	}

	if instance == nil {
		instance = new(Instance)
	}

	instance.Connections = connections
}

// NewConnections 实例化新的连接
func NewConnections(conf Conf) []*kafkaLib.Conn {
	var connections = make([]*kafkaLib.Conn, 0, len(conf.Endpoints))
	for _, addr := range conf.Endpoints {
		var (
			conn *kafkaLib.Conn
			err  error
		)
		//无账号密码连接
		if len(conf.Username) == 0 && len(conf.Password) == 0 {
			conn, err = kafkaLib.Dial("tcp", addr)
		} else {
			mechanism, mErr := getMechanism(conf)
			if mErr != nil {
				panic(mErr)
			}
			if conf.Timeout < 1 {
				conf.Timeout = 10000
			}
			dialer := &kafkaLib.Dialer{
				Timeout:       time.Duration(conf.Timeout) * time.Millisecond,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
			conn, err = dialer.DialContext(context.Background(), "tcp", addr)
		}
		if err != nil {
			panic(err)
		}
		connections = append(connections, conn)
	}

	return connections
}

// InitWriter 初始化写实例
func InitWriter(conf Conf, topic string) {
	writer := &kafkaLib.Writer{
		Addr:     kafkaLib.TCP(conf.Endpoints...),
		Topic:    topic,
		Balancer: &kafkaLib.Murmur2Balancer{},
	}

	var transport *kafkaLib.Transport
	if conf.Tls != nil {
		transport = &kafkaLib.Transport{TLS: conf.Tls}
	}
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		mechanism, mErr := getMechanism(conf)
		if mErr != nil {
			panic(mErr)
		}
		transport.SASL = mechanism
	}

	writer.Transport = transport

	if instance == nil {
		instance = new(Instance)
	}

	instance.Writer = writer
}

// NewWriter 实例化新的写实例
func NewWriter(conf Conf, topic string) (*kafkaLib.Writer, error) {
	writer := &kafkaLib.Writer{
		Addr:     kafkaLib.TCP(conf.Endpoints...),
		Topic:    topic,
		Balancer: &kafkaLib.Murmur2Balancer{},
	}

	var transport *kafkaLib.Transport
	if conf.Tls != nil {
		transport = &kafkaLib.Transport{TLS: conf.Tls}
	}
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		mechanism, mErr := getMechanism(conf)
		if mErr != nil {
			return nil, mErr
		}
		transport.SASL = mechanism
	}

	writer.Transport = transport

	return writer, nil
}

// InitReader 初始化读实例
func InitReader(conf Conf, topic, groupID string) {
	if conf.Timeout < 1 {
		conf.Timeout = 10000
	}
	dialer := &kafkaLib.Dialer{
		Timeout:   time.Duration(conf.Timeout) * time.Millisecond,
		DualStack: true,
	}

	if conf.Tls != nil {
		dialer.TLS = conf.Tls
	}

	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		mechanism, mErr := getMechanism(conf)
		if mErr != nil {
			panic(mErr)
		}
		dialer.SASLMechanism = mechanism
	}

	reader := kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers: conf.Endpoints,
		GroupID: groupID,
		Topic:   topic,
		Dialer:  dialer,
	})

	if instance == nil {
		instance = new(Instance)
	}

	instance.Reader = reader
}

// NewReader 实例化新的读实例
func NewReader(conf Conf, topic, groupID string) (*kafkaLib.Reader, error) {
	if conf.Timeout < 1 {
		conf.Timeout = 10000
	}
	dialer := &kafkaLib.Dialer{
		Timeout:   time.Duration(conf.Timeout) * time.Millisecond,
		DualStack: true,
	}

	if conf.Tls != nil {
		dialer.TLS = conf.Tls
	}

	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		mechanism, mErr := getMechanism(conf)
		if mErr != nil {
			return nil, mErr
		}
		dialer.SASLMechanism = mechanism
	}

	reader := kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers: conf.Endpoints,
		GroupID: groupID,
		Topic:   topic,
		Dialer:  dialer,
	})

	return reader, nil
}

// 根据SASL授权类型获取认证装置
func getMechanism(conf Conf) (sasl.Mechanism, error) {
	switch conf.SASLAuthType {
	default:
		return &plain.Mechanism{
			Username: conf.Username,
			Password: conf.Password,
		}, nil
	case "sha256":
		return scram.Mechanism(scram.SHA256, conf.Username, conf.Password)
	case "sha512":
		return scram.Mechanism(scram.SHA512, conf.Username, conf.Password)
	case "plain":
		return &plain.Mechanism{
			Username: conf.Username,
			Password: conf.Password,
		}, nil
	}
}
