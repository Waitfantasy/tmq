package rpc

import (
	"context"
	"github.com/Waitfantasy/tmq/message/manager"
	"github.com/Waitfantasy/tmq/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type MqServer struct {
	cfg     *Config
	manager *manager.Manager
}

func New(cfg *Config, mg *manager.Manager) *MqServer {
	return &MqServer{
		manager: mg,
		cfg:     cfg,
	}
}

func (s *MqServer) Prepare(ctx context.Context, request *pb.PrepareRequest) (*pb.PrepareResponse, error) {
	msg, err := s.manager.Prepare(request.Topic, int(request.RetrySecond), request.Body)
	if err != nil {
		return nil, err
	}

	return &pb.PrepareResponse{
		Id: msg.Id,
	}, nil
}

func (s *MqServer) Commit(ctx context.Context, request *pb.CommitRequest) (*pb.CommitResponse, error) {
	if _, err := s.manager.Send(request.Id); err != nil {
		return nil, err
	}

	return &pb.CommitResponse{
		Success: true,
	}, nil
}

func (s *MqServer) Rollback(ctx context.Context, request *pb.RollbackRequest) (*pb.RollbackReponse, error) {
	if _, err := s.manager.Cancel(request.Id); err != nil {
		return nil, err
	}

	return &pb.RollbackReponse{}, nil
}

func (s *MqServer) ConsumerNotify(ctx context.Context, ack *pb.ConsumerAck) (*pb.TMQVoid, error) {
	if _, err := s.manager.ConsumerCommit(ack.Id); err != nil {
		return nil, err
	}

	return &pb.TMQVoid{}, nil
}

func (s *MqServer) Run() error {
	var grpcServer *grpc.Server
	// create grpc server
	if s.cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFile, s.cfg.KeyFile)
		if err != nil {
			return err
		} else {
			grpcServer = grpc.NewServer(grpc.Creds(creds))
		}
	} else {
		grpcServer = grpc.NewServer()
	}

	// create listen
	if l, err := net.Listen("tcp", s.cfg.Addr); err != nil {
		return err
	} else {
		pb.RegisterTMQServiceServer(grpcServer, s)
		return grpcServer.Serve(l)
	}
}
