// Package main - точка входа приложения user_service.
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/obsessed-gopher/micro-service-guide/internal/adapters/hasher"
	"github.com/obsessed-gopher/micro-service-guide/internal/adapters/idgen"
	"github.com/obsessed-gopher/micro-service-guide/internal/adapters/repository"
	userservice "github.com/obsessed-gopher/micro-service-guide/internal/app/grpc/user_service"
	"github.com/obsessed-gopher/micro-service-guide/internal/config"
	"github.com/obsessed-gopher/micro-service-guide/internal/usecases"
	pb "github.com/obsessed-gopher/micro-service-guide/pkg/pb/user_service"
)

func main() {
	configPath := flag.String("config", "config/local.yml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Starting %s in %s mode", cfg.App.Name, cfg.App.Env)

	// Инициализация зависимостей
	// В реальном проекте здесь будет подключение к БД:
	// db, err := sql.Open("postgres", cfg.Database.DSN())
	// userRepo := repository.NewPostgresUserRepository(db)

	userRepo := repository.NewMemoryUserRepository()
	passwordHasher := hasher.NewBcryptHasher(0)
	idGenerator := idgen.NewUUIDGenerator()

	// Бизнес-логика
	userUsecase := usecases.NewUserUsecase(userRepo, passwordHasher, idGenerator)

	// gRPC сервер
	server := userservice.NewServer(userUsecase)

	// gRPC сервер
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, server)

	if cfg.App.Debug {
		reflection.Register(grpcServer)
	}

	listener, err := net.Listen("tcp", cfg.Server.GRPC.Addr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("gRPC server listening on %s", cfg.Server.GRPC.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
