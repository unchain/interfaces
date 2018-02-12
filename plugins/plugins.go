package plugins

import "github.com/unchainio/interfaces/logger"

type Message struct {
	Tag        uint64
	Body       []byte
	Attributes map[string]string
}

type EndpointPlugin interface {
	Init(config []byte, log logger.Logger) (err error)
	Send(message *Message) (response *Message, err error)
	Receive() (message *Message, err error)
	Ack(tag uint64, response *Message) error
	Nack(tag uint64) error
	Close() error
}

type ActionPlugin interface {
	Init(config []byte, log logger.Logger) (err error)
	Handle(message *Message) (result *Message, err error)
}
