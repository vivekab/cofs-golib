package golibgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"go.elastic.co/apm/module/apmgrpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxMsgSize = 30 << 20 //30MB
)

// Server is an interface describing the grpc server
type Server interface {
	GetGRPCServer() *ggrpc.Server
	Start(ctx context.Context, options ServerOptions) error
	Stop()
}

type ServerOptions struct {
	Port                 string
	Env                  string
	LocalRunningServices string
}

type server struct {
	server *ggrpc.Server
	sn     golibtypes.ServiceName
}

// NewServer returns a new server with middelware interceptors registered
func NewServer(ctx context.Context, serviceName golibtypes.ServiceName) Server {
	s := ggrpc.NewServer(
		ggrpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			apmgrpc.NewUnaryServerInterceptor(),
			unaryServerInterceptor(serviceName),
		)),
		ggrpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: 100 * time.Second,
		}),
		ggrpc.MaxSendMsgSize(maxMsgSize),
		ggrpc.MaxRecvMsgSize(maxMsgSize),
	)
	return &server{
		s,
		serviceName,
	}
}

func (s server) GetGRPCServer() *ggrpc.Server {
	return s.server
}

func (s server) Start(ctx context.Context, options ServerOptions) error {
	if options.Env == golibconstants.EnvLocal && isServiceRunningLocally(s.sn, options.LocalRunningServices) {
		options.Port = localServicePortMap[s.sn]
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", options.Port))
	if err != nil {
		return err
	}
	return s.server.Serve(l)
}

func (s server) Stop() {
	s.server.GracefulStop()
}
