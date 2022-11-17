package realmmgrgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type HealthCheckService struct {
	grpc_health_v1.UnimplementedHealthServer
	logger logging.Logger
	status grpc_health_v1.HealthCheckResponse_ServingStatus
}

func NewHealthChecker(logger logging.Logger) *HealthCheckService {
	return &HealthCheckService{
		logger: logger,
		status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
}

// Register health checker gRPC server implementation
func (s *HealthCheckService) Register(server *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(server, s)
}

func (s *HealthCheckService) Check(_ context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.logger.Debug("serving the check request for health check")
	return &grpc_health_v1.HealthCheckResponse{
		Status: s.status,
	}, nil
}

func (s *HealthCheckService) Watch(_ *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	s.logger.Debug("serving the watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: s.status,
	})
}
