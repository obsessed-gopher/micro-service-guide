package user_service

import (
	"context"

	pb "github.com/example/user-service/pkg/pb/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser возвращает пользователя по ID.
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.userModule.GetByID(ctx, req.Id)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetUserResponse{
		User: userToProto(user),
	}, nil
}