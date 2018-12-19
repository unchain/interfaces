package adapter

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCTriggerClient struct {
	broker *plugin.GRPCBroker
	client proto.TriggerClient
}

func (m *GRPCTriggerClient) Init(stub Stub, cfg []byte) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	_ = closer
	//defer closer()

	_, err := m.client.Init(context.Background(), &proto.InitTriggerRequest{
		StubServer: brokerID,
		Config:     cfg,
	})

	return err
}

func (m *GRPCTriggerClient) Trigger() (string, []byte, error) {
	r, err := m.client.Trigger(context.Background(), &proto.TriggerRequest{})

	if err != nil {
		return "", nil, err
	}

	return r.Tag, r.Message, nil
}

func (m *GRPCTriggerClient) Ack(tag string, response []byte) error {
	_, err := m.client.Ack(context.Background(), &proto.AckRequest{
		Tag:      tag,
		Response: response,
	})

	return err
}

func (m *GRPCTriggerClient) Nack(tag string, responseError error) error {
	_, err := m.client.Nack(context.Background(), &proto.NackRequest{
		Tag:   tag,
		Error: responseError.Error(),
	})

	return err
}

func (m *GRPCTriggerClient) Close() error {
	_, err := m.client.Close(context.Background(), &proto.CloseRequest{})

	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCTriggerServer struct {
	// This is the real implementation
	Impl   Trigger
	broker *plugin.GRPCBroker
}

func (m *GRPCTriggerServer) Init(ctx context.Context, req *proto.InitTriggerRequest) (*proto.InitTriggerResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	_ = closer
	//defer closer()

	return &proto.InitTriggerResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCTriggerServer) Trigger(ctx context.Context, req *proto.TriggerRequest) (*proto.TriggerResponse, error) {
	tag, r, err := m.Impl.Trigger()

	if err != nil {
		return nil, err
	}

	rBytes, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return &proto.TriggerResponse{
		Tag:     tag,
		Message: rBytes,
	}, nil
}

func (m *GRPCTriggerServer) Ack(ctx context.Context, req *proto.AckRequest) (*proto.AckResponse, error) {
	response := make(map[string]interface{})
	err := json.Unmarshal(req.Response, response)

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &proto.AckResponse{}, m.Impl.Ack(req.Tag, response)
}

func (m *GRPCTriggerServer) Nack(ctx context.Context, req *proto.NackRequest) (*proto.NackResponse, error) {
	return &proto.NackResponse{}, m.Impl.Nack(req.Tag, errors.New(req.Error))
}

func (m *GRPCTriggerServer) Close(ctx context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	return &proto.CloseResponse{}, m.Impl.Close()
}
