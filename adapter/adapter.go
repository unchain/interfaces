package adapter

import (
	"math/rand"

	"time"

	"github.com/unchainio/interfaces/logger"
)

type TaggedMessage struct {
	Tag uint64
	*Message
}

type Message struct {
	Body       []byte
	Attributes map[string]bool
}

type MessageOpts struct {
	tag uint64
}

func NewMessage(body []byte) *Message {
	return &Message{
		Body:       body,
		Attributes: make(map[string]bool),
	}
}

var defaultOpts = MessageOpts{}

type MessageOptsFunc func(opt *MessageOpts)

// NewMessage constructs a new message with a random tag unless a custom one is specified via WithTag(tag uint64)
func NewTaggedMessage(body []byte, optFuncs ...MessageOptsFunc) *TaggedMessage {
	opts := defaultOpts

	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}

	if opts.tag == 0 {
		opts.tag = randomTag()
	}

	return &TaggedMessage{
		Tag:     opts.tag,
		Message: NewMessage(body),
	}
}

func randomTag() uint64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
}

func WithTag(tag uint64) MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.tag = tag
	}
}

type Endpoint interface {
	Init(config []byte, log logger.Logger) (err error)
	Send(message *Message) (response *Message, err error)
	Receive() (message *TaggedMessage, err error)
	Ack(tag uint64, response *Message) error
	Nack(tag uint64) error
	Close() error
}

type Action interface {
	Init(config []byte, log logger.Logger) (err error)
	Invoke(message *Message) (err error)
}
