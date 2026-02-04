package user_service

import (
	"context"

	"github.com/obsessed-gopher/micro-service-guide/internal/types"
	"github.com/obsessed-gopher/micro-service-guide/internal/usecases"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
)

// ListUsers возвращает список пользователей.
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	filter := usecases.ListFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	if req.Filter != nil {
		filter.IDs = req.Filter.Ids
		filter.Emails = req.Filter.Emails
		filter.Statuses = statusesFromProto(req.Filter.Statuses)
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

func statusesFromProto(statuses []pb.UserStatus) []types.UserStatus {
	if len(statuses) == 0 {
		return nil
	}

	result := make([]types.UserStatus, len(statuses))
	for i, s := range statuses {
		result[i] = statusFromProto(s)
	}

	return result
}
