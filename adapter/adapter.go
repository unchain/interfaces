package adapter

import "github.com/unchainio/interfaces/logger"

type ComponentType string

const (
	TriggerComponent = "trigger"
	ActionComponent  = "action"
)

type Trigger interface {
	// Init: must NOT block, start long running processes in a go routine
	Init(stub Stub, config []byte) (err error)

	// Trigger: must block until a new message is received
	Trigger() (tag string, message map[string]interface{}, err error)

	// Ack is called by the adapter base after the message (with tag `tag`), which was initially received
	// by this input endpoint, has been successfully passed through the actions in the pipeline, sent
	// over the output endpoint, and a response has been returned from the output endpoint and passed
	// through the actions in the response pipeline.
	Ack(tag string, response map[string]map[string]interface{}) error

	// Nack is called by the adapter base if anything goes wrong while processing the message with tag `tag`
	Nack(tag string, err error) error

	Close() error
}

type Action interface {
	Init(stub Stub, config []byte) (err error)
	Invoke(inputMessage map[string]map[string]interface{}) (outputMessage map[string]interface{}, err error)
}

type Stub interface {
	logger.Logger

	// TODO in the future this interface will also contain a kv store and a secret store
}
