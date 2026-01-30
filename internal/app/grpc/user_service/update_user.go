package user_service

import (
	"context"

	"github.com/example/user-service/internal/models"
	pb "github.com/example/user-service/pkg/pb/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser обновляет данные пользователя.
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	input := models.UpdateUserInput{}
	if req.Email != nil {
		input.Email = req.Email
	}
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Status != nil {
		st := statusFromProto(*req.Status)
		input.Status = &st
	}

	user, err := s.userModule.Update(ctx, req.Id, input)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.UpdateUserResponse{
		User: userToProto(user),
	}, nil
}