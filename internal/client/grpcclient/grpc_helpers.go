package grpcclient

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func (c *GRPCClient) ConnectionNumber() uint64 {

	return c.connectionNumber
}

func (c *GRPCClient) withAuthRetry(ctx context.Context, fn func(context.Context) error) error {
	if err := c.ensureAuth(ctx); err != nil {
		return err
	}

	authCtx := c.createAuthContext(ctx)
	err := fn(authCtx)

	if isUnauthorizedError(err) {
		c.token = ""
		if err := c.ensureAuth(ctx); err != nil {
			return err
		}

		authCtx = c.createAuthContext(ctx)
		return fn(authCtx)
	}

	return err
}

func (c *GRPCClient) createAuthContext(ctx context.Context) context.Context {

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.NewOutgoingContext(
			ctx,
			metadata.Pairs("authorization", c.token),
		)
	}

	md.Append("authorization", c.token)
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *GRPCClient) ensureAuth(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	return c.Login(ctx, c.login, c.password)
}

func isUnauthorizedError(err error) bool {
	return err != nil && err.Error() == "rpc error: code = Unauthenticated desc = invalid token"
}

func (c *GRPCClient) withConnNumber(ctx context.Context) context.Context {

	return metadata.NewOutgoingContext(ctx,
		metadata.Pairs("connectionnumber", strconv.FormatUint(c.connectionNumber, 10)),
	)
}
