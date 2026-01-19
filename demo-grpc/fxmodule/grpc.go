package fxmodule

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	localConfig "github.com/gohive/demo-grpc/config"
)

func NewGRPCServer(lc fx.Lifecycle) *grpc.Server {
	cfg := localConfig.Cfg.GRPC

	// Create gRPC server with interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			loggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)

	// Enable reflection for grpcurl/grpcui in development
	if config.IsDevelopment() {
		reflection.Register(server)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			go func() {
				logger.Infof("gRPC server starting on port %d", cfg.Port)
				if err := server.Serve(lis); err != nil {
					logger.Errorf("gRPC server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("gRPC server shutting down")
			server.GracefulStop()
			return nil
		},
	})

	return server
}

// recoveryInterceptor handles panics
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic recovered in gRPC: %v", r)
			err = fmt.Errorf("internal server error")
		}
	}()
	return handler(ctx, req)
}

// loggingInterceptor logs requests
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Debugf("gRPC request: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Errorf("gRPC error: %s - %v", info.FullMethod, err)
	}
	return resp, err
}
