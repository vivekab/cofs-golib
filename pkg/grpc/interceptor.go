package golibgrpc

import (
	"context"
	"runtime/debug"
	"time"

	goliberror "github.com/vivekab/golib/pkg/error"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibcontext "github.com/vivekab/golib/pkg/context"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"github.com/vivekab/golib/protobuf/protoroot"

	golibid "github.com/vivekab/golib/pkg/id"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	impl "google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"

	goliblogging "github.com/vivekab/golib/pkg/logging"

	"github.com/mohae/deepcopy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	requestIDKey = ctxKey("requestID")
	loggerKey    = ctxKey("rLogger")
)

type ctxKey string

func MaskSensitive(data interface{}) interface{} {
	if data == nil {
		return data
	}
	proto.Message(data.(proto.Message)).ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		for _, extFlag := range []*impl.ExtensionInfo{protoroot.E_Sensitive, protoroot.E_Elongated} {
			s := proto.GetExtension(opts, extFlag)
			if flag, ok := s.(bool); ok && flag {
				data.(proto.Message).ProtoReflect().Clear(fd)
				return true
			}
		}
		return true
	})
	return data
}

func unaryServerInterceptor(serviceName golibtypes.ServiceName) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		var h interface{}

		start := time.Now()

		ri := golibcontext.GetFromContext(ctx, golibconstants.HeaderRequestID)
		requestCopy := deepcopy.Copy(req)
		fields := goliblogging.Fields{
			"method":  info.FullMethod,
			"request": MaskSensitive(requestCopy),
		}
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "panic: %v", r)
				fields["error"] = err.Error()
				fields["stack"] = goliblogging.PrettifyStack(string(debug.Stack()))
			}
			if info.FullMethod != "/grpc.health.v1.Health/Check" {
				fields["duration"] = time.Since(start) / time.Millisecond
				goliblogging.InfoD(ctx, "GRPC Call LOG", fields)
			}
		}()
		ctx = context.WithValue(ctx, requestIDKey, ri)
		ctx = context.WithValue(ctx, loggerKey, goliblogging.GetLogger())

		h, err = handler(ctx, req)
		hCopy := deepcopy.Copy(h)
		fields["response"] = MaskSensitive(hCopy)
		if err != nil {
			fields["err"] = goliberror.GetGrpcError(err)
		}
		return h, err
	}
}

func unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := make(metadata.MD)
		if omd, ok := metadata.FromOutgoingContext(ctx); ok {
			md = metadata.Join(md, omd)
		}
		if imd, ok := metadata.FromIncomingContext(ctx); ok {
			md = metadata.Join(md, imd)
		}
		if _, ok := md[golibconstants.HeaderRequestID]; !ok {
			requestId, _ := golibid.NewId(golibid.IdPrefixRequest)
			md[golibconstants.HeaderRequestID] = append(md[golibconstants.HeaderRequestID], requestId.String())
		}
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
	}
}
