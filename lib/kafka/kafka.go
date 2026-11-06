package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/segmentio/kafka-go/sasl/scram"

	kafkaLib "github.com/segmentio/kafka-go"
)

// Instance connection instance
type Instance struct {
	Connections []*kafkaLib.Conn
	Writer      *kafkaLib.Writer
	Reader      *kafkaLib.Reader
}

var instance *Instance

// GetConnections retrieves connections
//
// Call Init to initialize before calling this method
func GetConnections() []*kafkaLib.Conn {
	return instance.Connections
}

// GetWriter retrieves the writer instance
//
// Call Init to initialize before calling this method
func GetWriter() *kafkaLib.Writer {
	return instance.Writer
}

// GetReader retrieves the reader instance
//
// Call Init to initialize before calling this method
func GetReader() *kafkaLib.Reader {
	return instance.Reader
}

// GetInstance retrieves the complete instance
//
// Includes connections, reader, and writer instances
func GetInstance() *Instance {
	return instance
}

// Init initializes connections
//
// This method initializes connections, reader, and writer instances
//
// After initialization, call GetInstance to retrieve the complete instance
func Init(conf Conf, topic, groupID string) {
	InitConnections(conf)
	InitWriter(conf, topic)
	InitReader(conf, topic, groupID)
}

// New initializes connections
//
// This method initializes connections, reader, and writer instances
func New(conf Conf, topic, groupID string) (connections []*kafkaLib.Conn,
	writer *kafkaLib.Writer, wErr error,
	reader *kafkaLib.Reader, rErr error) {
	connections = NewConnections(conf)
	writer, wErr = NewWriter(conf, topic)
	reader, rErr = NewReader(conf, topic, groupID)

	return
}

// InitConnections initializes connections
func InitConnections(conf Conf) {
	var connections = make([]*kafkaLib.Conn, 0, len(conf.Endpoints))
	for _, addr := range conf.Endpoints {
		var (
			conn *kafkaLib.Conn
			err  error
		)
		//Connect without a username or password
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

// NewConnections instantiates new connections
func NewConnections(conf Conf) []*kafkaLib.Conn {
	var connections = make([]*kafkaLib.Conn, 0, len(conf.Endpoints))
	for _, addr := range conf.Endpoints {
		var (
			conn *kafkaLib.Conn
			err  error
		)
		//Connect without a username or password
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

// InitWriter initializes the writer instance
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

// NewWriter instantiates a new writer
func NewWriter(conf Conf, topic string) (*kafkaLib.Writer, error) {
	writer := &kafkaLib.Writer{
		Addr:     kafkaLib.TCP(conf.Endpoints...),
		Topic:    topic,
		Balancer: &kafkaLib.Murmur2Balancer{},
	}

	var transport = &kafkaLib.Transport{}
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

// InitReader initializes the reader instance
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

// NewReader instantiates a new reader
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

// Get the authentication mechanism for the SASL authentication type
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
