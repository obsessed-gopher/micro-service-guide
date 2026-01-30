package user_service

import (
	"context"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser создаёт нового пользователя.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := s.userModule.Create(ctx, models.CreateUserInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.CreateUserResponse{
		User: userToProto(user),
	}, nil
}
