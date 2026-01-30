// Package user_service содержит gRPC сервер и хендлеры.
package user_service

import (
	"context"

	"github.com/example/user-service/internal/models"
	pb "github.com/example/user-service/pkg/pb/user_service"
)

// UserModule - интерфейс бизнес-логики пользователей.
type UserModule interface {
	Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error)
}

// Server - gRPC сервер сервиса пользователей.
type Server struct {
	pb.UnimplementedUserServiceServer
	userModule UserModule
}

// NewServer создаёт новый сервер.
func NewServer(userModule UserModule) *Server {
	return &Server{
		userModule: userModule,
	}
}