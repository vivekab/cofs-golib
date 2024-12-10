package golibgrpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"go.elastic.co/apm/module/apmgrpc"
	ggrpc "google.golang.org/grpc"
)

const timeoutSec = 100

// Client encapsulates grpc client details
type Client interface {
	GetConn() *ggrpc.ClientConn
	GetContext() context.Context
	CloseAndCancel()
}

type ClientOptions struct {
	Name                 golibtypes.ServiceName
	Env                  string // TODO - add comments
	ConnectEnv           string
	LocalRunningServices string
	Port                 string
}

type client struct {
	conn   *ggrpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func (c client) GetConn() *ggrpc.ClientConn {
	return c.conn
}

func (c client) GetContext() context.Context {
	return c.ctx
}

func (c client) CloseAndCancel() {
	c.cancel()
	c.conn.Close()
}

func NewClient(ctx context.Context, options ClientOptions) (Client, error) {
	var err error
	var conn *ggrpc.ClientConn
	var connStr string

	ctx, cancel := context.WithTimeout(ctx, timeoutSec*time.Second)
	defer cancel() // Ensure cancel is called before returning

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable, codes.Aborted),
		grpc_retry.WithMax(3),
	}
	// Get connection string for the service
	if options.Port == "" {
		options.Port = "3001"
	}
	if connStr, err = GetConnectionStringForService(ctx, options); err != nil {
		return nil, err
	}

	// Dial the gRPC connection
	conn, err = ggrpc.NewClient(connStr, ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		ggrpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(opts...),
			unaryClientInterceptor(),
			apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		return nil, err
	}

	// Return the client with the connection and context
	return &client{
		conn,
		ctx,
		cancel, // Pass cancel function to the client
	}, nil
}
