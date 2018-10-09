package adapter

import (
	"math/rand"
	"time"

	"github.com/unchainio/pkg/xsync"
)

var globalCounter xsync.Counter

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ITaggedMessage interface {
	GetTag() uint64

	GetID() string
	GetBody() []byte
}

type IMessage interface {
	GetID() string
	GetBody() []byte
}


type TaggedMessage struct {
	Tag uint64
	*Message
}

func (tm *TaggedMessage) GetTag() uint64 {
	return tm.Tag
}

type Message struct {
	ID         string
	Body      []byte
	Attributes map[string]bool
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) GetBody() []byte {
	return m.Body
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
