package user_service

import (
	"context"

	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUser удаляет пользователя.
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.userModule.Delete(ctx, req.Id); err != nil {
		return nil, mapError(err)
	}

	return &pb.DeleteUserResponse{}, nil
}
