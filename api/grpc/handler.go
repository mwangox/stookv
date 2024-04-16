package grpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"stoo-kv/api"
	"stoo-kv/api/grpc/proto"
	"stoo-kv/config"
	"stoo-kv/internal/crypto"
	"stoo-kv/internal/store"
)

type Server struct {
	storage store.Store
	config  *config.Config
	proto.UnimplementedKVServiceServer
}

func NewGrpcServer(storage store.Store, config *config.Config) *Server {
	return &Server{
		config:  config,
		storage: storage,
	}
}
func (s *Server) GetService(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	value, err := s.storage.Get(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key))
	if err != nil {
		message := fmt.Sprintf("Failed to read keys from storage: %v", err)
		log.Printf(message)
		return nil, status.Errorf(codes.Aborted, message)
	}

	if value == "" {
		message := "key not found from storage"
		log.Printf(message)
		return nil, status.Errorf(codes.NotFound, message)
	}
	value, err = api.CheckEncryption(value, s.config)
	if err != nil {
		log.Printf("Failed to decrypt the value: %v", err)
		return nil, status.Errorf(codes.Aborted, "data decryption failed")
	}
	return &proto.GetResponse{Data: value}, nil
}

//func (s *Server) GetAllService(ctx context.Context, request *proto.GetAllRequest) (*proto.GetAllResponse, error) {
//	values, err := s.storage.GetAll()
//	if err != nil {
//		log.Printf("Failed to read keys from storage: %v", err)
//		return nil, err
//	}
//	if len(values) == 0 {
//		log.Printf("Keys not found from storage")
//		return nil, errors.New("keys not found from storage")
//	}
//	return &proto.GetAllResponse{Data: api.ParseValues(values, s.config)}, nil
//}

func (s *Server) GetServiceByNamespaceAndProfile(ctx context.Context, request *proto.GetByNamespaceAndProfileRequest) (*proto.GetByNamespaceAndProfileResponse, error) {
	values, err := s.storage.GetByNameSpaceAndProfile(request.Namespace, request.Profile)
	if err != nil {
		message := fmt.Sprintf("Failed to read keys from storage: %v", err)
		log.Printf(message)
		return nil, status.Errorf(codes.Aborted, message)
	}
	if len(values) == 0 {
		message := "keys not found from storage"
		log.Printf(message)
		return nil, status.Errorf(codes.NotFound, message)
	}
	return &proto.GetByNamespaceAndProfileResponse{Data: api.ParseValues(values, s.config)}, nil
}

func (s *Server) SetKeyService(ctx context.Context, request *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	return s.set(request, false)
}

func (s *Server) SetSecretKeyService(ctx context.Context, request *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	return s.set(request, true)
}

func (s *Server) set(request *proto.SetKeyRequest, isSecret bool) (*proto.SetKeyResponse, error) {
	value := request.Value
	if isSecret {
		ciphertext, err := crypto.Encrypt([]byte(value), s.config.Application.EncryptKey)
		if err != nil {
			log.Printf("Failed to encrypt data: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		encPrefix := s.config.Application.EncryptPrefix
		if encPrefix == "" {
			encPrefix = "{ENC} "
		}

		value = encPrefix + hex.EncodeToString(ciphertext)
	}
	if err := s.storage.Set(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key), value); err != nil {
		log.Printf("Failed to store data into storage: %v", err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &proto.SetKeyResponse{Data: "Data saved successfully"}, nil
}

func (s *Server) DeleteKeyService(ctx context.Context, request *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	if err := s.storage.Delete(fmt.Sprintf("%s::%s::%s", request.Namespace, request.Profile, request.Key)); err != nil {
		log.Printf("Failed to remove data from storage: %v", err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &proto.DeleteKeyResponse{Data: "Key removed successfully"}, nil
}

func RunGrpcServer(cfg *config.Config, storage store.Store) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Application.GrpcPort))
	if err != nil {
		return err
	}
	var options []grpc.ServerOption
	if cfg.Application.GrpcUseTls {
		creds, err := credentials.NewServerTLSFromFile(cfg.Application.GrpcServerCert, cfg.Application.GrpcServerKey)
		if err != nil {
			log.Fatalf("Failed to create grpc credentials: %v", err)
		}
		options = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s := grpc.NewServer(options...)
	reflection.Register(s)
	proto.RegisterKVServiceServer(s, NewGrpcServer(storage, cfg))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start grpc server: %v", err)
		}
	}()
	return nil
}
