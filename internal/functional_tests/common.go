package functional_tests

import (
	"context"
	"google.golang.org/grpc/test/bufconn"
	"net"
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
