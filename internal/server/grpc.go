package server

import (
	v1 "gitlab.calendaria.team/services/stores/api/stores/v1"
	"gitlab.calendaria.team/services/stores/internal/conf"
	"gitlab.calendaria.team/services/stores/internal/service"
	u_jwt "gitlab.calendaria.team/services/utils/v4/jwt"
	u_auth "gitlab.calendaria.team/services/utils/v4/middlewares/auth"
	u_tracing "gitlab.calendaria.team/services/utils/v4/tracing"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Bootstrap,
	jwtp u_jwt.IJwtProcessor,
	tracer u_tracing.ITracer,
	storesService *service.StoresService,
) *grpc.Server {
	err := tracer.Initialize()
	if err != nil {
		panic(err)
	}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			tracing.Server(),
			u_auth.Server(jwtp),
		),
	}
	if c.GetServer().GetGrpc().GetNetwork() != "" {
		opts = append(opts, grpc.Network(c.GetServer().GetGrpc().GetNetwork()))
	}
	if c.GetServer().GetGrpc().GetAddr() != "" {
		opts = append(opts, grpc.Address(c.GetServer().GetGrpc().GetAddr()))
	}
	if c.GetServer().GetGrpc().GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(c.GetServer().GetGrpc().GetTimeout().AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	v1.RegisterStoresServer(srv, storesService)

	return srv
}
