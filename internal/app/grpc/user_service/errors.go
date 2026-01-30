package user_service

import (
	"github.com/example/user-service/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapError конвертирует бизнес-ошибки в gRPC статусы.
func mapError(err error) error {
	switch err {
	case types.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case types.ErrUserAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case types.ErrInvalidEmail, types.ErrInvalidPassword:
		return status.Error(codes.InvalidArgument, err.Error())
	case types.ErrUserBlocked:
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}