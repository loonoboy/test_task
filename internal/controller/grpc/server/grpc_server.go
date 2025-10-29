package server

import (
	"context"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller"
	"git.amocrm.ru/study_group/in_memory_database/internal/controller/grpc/proto"
)

type GRPCServer struct {
	proto.UnimplementedDeleteAccountServiceServer
	AccountService controller.AccountUsecaseInterface
}

func NewGRPCServer(accountService controller.AccountUsecaseInterface) *GRPCServer {
	return &GRPCServer{
		AccountService: accountService,
	}
}

func (s *GRPCServer) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	err := s.AccountService.DeleteAccount(int(req.AccountId))

	if err != nil {
		return &proto.DeleteAccountResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &proto.DeleteAccountResponse{
		Success: true,
		Message: "Success",
	}, err
}
