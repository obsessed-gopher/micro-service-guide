package user_service

import (
	"github.com/example/user-service/internal/models"
	"github.com/example/user-service/internal/types"
	pb "github.com/example/user-service/pkg/pb/user_service"
)

// userToProto конвертирует бизнес-модель в proto.
func userToProto(u *models.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Status:    statusToProto(u.Status),
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}

// statusToProto конвертирует внутренний статус в proto.
func statusToProto(s types.UserStatus) pb.UserStatus {
	switch s {
	case types.UserStatusActive:
		return pb.UserStatus_USER_STATUS_ACTIVE
	case types.UserStatusInactive:
		return pb.UserStatus_USER_STATUS_INACTIVE
	case types.UserStatusBlocked:
		return pb.UserStatus_USER_STATUS_BLOCKED
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// statusFromProto конвертирует proto статус во внутренний.
func statusFromProto(s pb.UserStatus) types.UserStatus {
	switch s {
	case pb.UserStatus_USER_STATUS_ACTIVE:
		return types.UserStatusActive
	case pb.UserStatus_USER_STATUS_INACTIVE:
		return types.UserStatusInactive
	case pb.UserStatus_USER_STATUS_BLOCKED:
		return types.UserStatusBlocked
	default:
		return types.UserStatusUnspecified
	}
}