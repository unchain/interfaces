package adapter

import (
	"math/rand"
	"time"

	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/xsync"
)

var globalCounter xsync.Counter

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
		opts.tag = globalCounter.Add(1)
	}

	return &TaggedMessage{
		Tag:     opts.tag,
		Message: NewMessage(body),
	}
}

func randomTag() uint64 {
	return rand.Uint64()
}

func WithTag(tag uint64) MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.tag = tag
	}
}

func WithRandomTag() MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.tag = randomTag()
	}
}

type Endpoint interface {
	// Init: must NOT block, start long running processes in a go routine
	Init(config []byte, log logger.Logger) (err error)

	// Send: must block until sending is complete
	Send(message *Message) (response *Message, err error)

	// Receive: must block until a new message is received
	Receive() (message *TaggedMessage, err error)

	// Ack is called by the adapter base after the message (with tag `tag`), which was initially received
	// by this input endpoint, has been successfully passed through the actions in the pipeline, sent
	// over the output endpoint, and a response has been returned from the output endpoint and passed
	// through the actions in the response pipeline.
	Ack(tag uint64, response *Message) error

	// Nack is called by the adapter base if anything goes wrong while processing the message with tag `tag`
	Nack(tag uint64) error
	Close() error
}

type Action interface {
	Init(config []byte, log logger.Logger) (err error)
	Invoke(message *Message) (err error)
}
