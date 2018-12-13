package adapter

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCEndpointClient struct {
	broker *plugin.GRPCBroker
	client proto.EndpointClient
}

func (m *GRPCEndpointClient) Init(stub Stub, cfg []byte) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	_ = closer
	//defer closer()

	_, err := m.client.Init(context.Background(), &proto.InitEndpointRequest{
		StubServer: brokerID,
		Config:     cfg,
	})

	return err
}

func (m *GRPCEndpointClient) Receive() (string, []byte, error) {
	r, err := m.client.Receive(context.Background(), &proto.ReceiveRequest{})

	if err != nil {
		return "", nil, err
	}

	return r.Tag, r.Message, nil
}

func (m *GRPCEndpointClient) Ack(tag string, response []byte) error {
	_, err := m.client.Ack(context.Background(), &proto.AckRequest{
		Tag:      tag,
		Response: response,
	})

	return err
}

func (m *GRPCEndpointClient) Nack(tag string, responseError error) error {
	_, err := m.client.Nack(context.Background(), &proto.NackRequest{
		Tag:   tag,
		Error: responseError.Error(),
	})

	return err
}

func (m *GRPCEndpointClient) Close() error {
	_, err := m.client.Close(context.Background(), &proto.CloseRequest{})

	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCEndpointServer struct {
	// This is the real implementation
	Impl   Endpoint
	broker *plugin.GRPCBroker
}

func (m *GRPCEndpointServer) Init(ctx context.Context, req *proto.InitEndpointRequest) (*proto.InitEndpointResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	_ = closer
	//defer closer()

	return &proto.InitEndpointResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCEndpointServer) Receive(ctx context.Context, req *proto.ReceiveRequest) (*proto.ReceiveResponse, error) {
	tag, r, err := m.Impl.Receive()

	if err != nil {
		return nil, err
	}

	rBytes, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return &proto.ReceiveResponse{
		Tag:     tag,
		Message: rBytes,
	}, nil
}

func (m *GRPCEndpointServer) Ack(ctx context.Context, req *proto.AckRequest) (*proto.AckResponse, error) {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(req.Response, response)

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &proto.AckResponse{}, m.Impl.Ack(req.Tag, response)
}

func (m *GRPCEndpointServer) Nack(ctx context.Context, req *proto.NackRequest) (*proto.NackResponse, error) {
	return &proto.NackResponse{}, m.Impl.Nack(req.Tag, errors.New(req.Error))
}

func (m *GRPCEndpointServer) Close(ctx context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	return &proto.CloseResponse{}, m.Impl.Close()
}
