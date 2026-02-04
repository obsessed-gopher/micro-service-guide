// Package user_service содержит gRPC сервер и хендлеры.
package user_service

import (
	"context"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
)

// UserUsecase - интерфейс бизнес-логики пользователей.
type UserUsecase interface {
	Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error)
}

// Server - gRPC сервер сервиса пользователей.
type Server struct {
	pb.UnimplementedUserServiceServer
	userUsecase UserUsecase
}

// NewServer создаёт новый сервер.
func NewServer(userUsecase UserUsecase) *Server {
	return &Server{
		userUsecase: userUsecase,
	}
}
