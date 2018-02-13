package plugins

import (
	"math/rand"

	"time"

	"github.com/unchainio/interfaces/logger"
)

type Message struct {
	Tag        uint64
	Body       []byte
	Attributes map[string]bool
}

type MessageOpts struct {
	Tag uint64
}

var defaultOpts = &MessageOpts{}

type MessageOptsFunc func(opt *MessageOpts)

// NewMessage constructs a new message with a random tag unless a custom one is specified via WithTag(tag uint64)
func NewMessage(body []byte, optFuncs ...MessageOptsFunc) *Message {
	opts := defaultOpts

	for _, optFunc := range optFuncs {
		optFunc(opts)
	}

	if opts.Tag == 0 {
		opts.Tag = randomTag()
	}

	return &Message{
		Tag:        opts.Tag,
		Body:       body,
		Attributes: make(map[string]bool),
	}
}

func randomTag() uint64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
}

func WithTag(tag uint64) MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.Tag = tag
	}
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
	Invoke(message *Message) (result *Message, err error)
}
