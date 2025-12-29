package functional_tests

import (
	"context"
	grpcserver "github.com/s-turchinskiy/keeper/internal/server/grpc"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/internal/server/token"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"time"
)

const (
	loginNewUser      = "new user"
	loginExistingUser = "existing user"
	password          = "password"

	secretOnlyCreating = "Secret only creating"
	secretForDeleting  = "Secret for deleting"
	secretForUpdating  = "Secret for updating"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func runGRPCServer(usersRepository repository.UserRepositorier, secretRepository repository.SecretRepositorier) {
	lis = bufconn.Listen(bufSize)

	jwtManager := token.NewJWTManager("secret", time.Minute)
	srvc := service.NewService(
		jwtManager,
		usersRepository,
		secretRepository,
	)
	grpcServer := grpcserver.NewGrpcServer(srvc)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}
