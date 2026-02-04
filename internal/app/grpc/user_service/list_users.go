package user_service

import (
	"context"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
)

// ListUsers возвращает список пользователей.
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	filter := models.ListUsersFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}
	if req.Status != nil {
		st := statusFromProto(*req.Status)
		filter.Status = &st
	}

	users, total, err := s.userUsecase.List(ctx, filter)
	if err != nil {
		return nil, mapError(err)
	}

	protoUsers := make([]*pb.User, len(users))
	for i, u := range users {
		protoUsers[i] = userToProto(u)
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
	}, nil
}
