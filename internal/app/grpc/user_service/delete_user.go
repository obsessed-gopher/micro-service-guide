package user_service

import (
	"context"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUser удаляет пользователя.
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	filter := models.UserFilter{IDs: []string{req.Id}}
	if _, err := s.userUsecase.Delete(ctx, filter); err != nil {
		return nil, mapError(err)
	}

	return &pb.DeleteUserResponse{}, nil
}
