package rpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"stoo-kv/api"
	"stoo-kv/api/rpc/proto"
	"stoo-kv/config"
	"stoo-kv/internal"
)

type GrpcServer struct {
	storage internal.Store
	config  *config.Config
	proto.UnimplementedKVServiceServer
}

func NewGrpcServer(storage internal.Store, config *config.Config) *GrpcServer {
	return &GrpcServer{
		config:  config,
		storage: storage,
	}
}
func (s *GrpcServer) GetService(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	value, err := s.storage.Get(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key))
	if err != nil {
		log.Printf("Failed to read key from storage: %v", err)
		return nil, err
	}

	if value == "" {
		return nil, errors.New("data not found from storage")
	}

	value, err = api.CheckEncryption(value, s.config)
	if err != nil {
		log.Printf("Failed to decrypt the value: %v", err)
		return nil, errors.New("data decryption failed")
	}
	return &proto.GetResponse{Data: value}, nil
}

func (s *GrpcServer) GetAllService(ctx context.Context, request *proto.GetAllRequest) (*proto.GetAllResponse, error) {
	values, err := s.storage.GetAll()
	if err != nil {
		log.Printf("Failed to read keys from storage: %v", err)
		return nil, err
	}
	if len(values) == 0 {
		log.Printf("Keys not found from storage")
		return nil, errors.New("keys not found from storage")
	}
	return &proto.GetAllResponse{Data: api.ParseValues(values, s.config)}, nil
}

func (s *GrpcServer) GetServiceByNamespaceAndProfile(ctx context.Context, request *proto.GetByNamespaceAndProfileRequest) (*proto.GetByNamespaceAndProfileResponse, error) {
	values, err := s.storage.GetByNameSpaceAndProfile(request.Namespace, request.Profile)
	if err != nil {
		log.Printf("Failed to read keys from storage: %v", err)
		return nil, err
	}
	if len(values) == 0 {
		log.Printf("Keys not found from storage")
		return nil, errors.New("keys not found from storage")
	}
	return &proto.GetByNamespaceAndProfileResponse{Data: api.ParseValues(values, s.config)}, nil
}

func (s *GrpcServer) SetKeyService(ctx context.Context, request *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	if err := s.storage.Set(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key), request.Value); err != nil {
		log.Printf("Failed to store data into storage: %v", err)
		return nil, err
	}
	return &proto.SetKeyResponse{Data: "Data saved successfully"}, nil
}

func (s *GrpcServer) DeleteKeyService(ctx context.Context, request *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	if err := s.storage.Delete(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key)); err != nil {
		log.Printf("Failed to remove data from storage: %v", err)
		return nil, err
	}
	return &proto.DeleteKeyResponse{Data: "Key removed successfully"}, nil
}

func RunGrpcServer(cfg *config.Config, storage internal.Store) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcPort))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	reflection.Register(s)
	proto.RegisterKVServiceServer(s, NewGrpcServer(storage, cfg))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start grpc server: %v", err)
		}
	}()
	return nil
}
