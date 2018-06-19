package adapter

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"github.com/unchainio/interfaces/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.ActionClient
	log    logger.Logger
}

func (m *GRPCClient) Init(cfg []byte, log logger.Logger) error {
	logHelperServer := &GRPCLogServer{Impl: log}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterLogHelperServer(s, logHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	_, err := m.client.Init(context.Background(), &proto.InitRequest{
		LogServer: brokerID,
		Config:    cfg,
	})

	s.Stop()
	return err
}

func (m *GRPCClient) Invoke(message *Message) error {
	imsg, err := m.client.Invoke(context.Background(), &proto.InvokeMessage{
		Body:       message.Body,
		Attributes: message.Attributes,
	})

	if err != nil {
		return err
	}

	message.Body = imsg.Body
	message.Attributes = imsg.Attributes

	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl   Action
	broker *plugin.GRPCBroker
}

func (m *GRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	conn, err := m.broker.Dial(req.LogServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	log := &GRPCLogHelperClient{proto.NewLogHelperClient(conn)}
	return &proto.Empty{}, m.Impl.Init(req.Config, log)
}

func (m *GRPCServer) Invoke(ctx context.Context, req *proto.InvokeMessage) (*proto.InvokeMessage, error) {
	msg := &Message{
		Body:       req.Body,
		Attributes: req.Attributes,
	}

	err := m.Impl.Invoke(msg)

	return &proto.InvokeMessage{
		Body:       msg.Body,
		Attributes: msg.Attributes,
	}, err
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCLogHelperClient struct{ client proto.LogHelperClient }

func (m *GRPCLogHelperClient) Debugf(format string, v ...interface{}) {
	m.client.Debugf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCLogHelperClient) Errorf(format string, v ...interface{}) {
	m.client.Errorf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCLogHelperClient) Fatalf(format string, v ...interface{}) {
	m.client.Fatalf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCLogHelperClient) Panicf(format string, v ...interface{}) {
	m.client.Panicf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCLogHelperClient) Printf(format string, v ...interface{}) {
	m.client.Printf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCLogHelperClient) Warnf(format string, v ...interface{}) {
	m.client.Warnf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCLogServer struct {
	// This is the real implementation
	Impl logger.Logger
}

func (m *GRPCLogServer) Printf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Printf("%s", req.Message)

	return &proto.Empty{}, nil
}

func (m *GRPCLogServer) Fatalf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Fatalf("%s", req.Message)

	return &proto.Empty{}, nil
}

func (m *GRPCLogServer) Panicf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Panicf("%s", req.Message)

	return &proto.Empty{}, nil
}

func (m *GRPCLogServer) Debugf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Debugf("%s", req.Message)

	return &proto.Empty{}, nil
}

func (m *GRPCLogServer) Warnf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Warnf("%s", req.Message)

	return &proto.Empty{}, nil
}

func (m *GRPCLogServer) Errorf(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	m.Impl.Errorf("%s", req.Message)

	return &proto.Empty{}, nil
}
