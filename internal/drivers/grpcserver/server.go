package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

const grpcPortStart = 1024

// Service represents a service that can register itself with the gRPC server.
type Service interface {
	Register(server *grpc.Server)
}

// Server is the gRPC server that can be run.
type Server struct {
	server *grpc.Server
	port   uint16
}

// NewServer constructs a new grpcServer, registering the provided services.
// Service should not be running as root, thus provided port should not be < 1024
func NewServer(port uint16, services []Service, opt ...grpc.ServerOption) (*Server, error) {
	if port < grpcPortStart {
		return nil, realmmgr_errors.NewInvalidArgumentError("port", "cannot be less than 1024")
	}
	server := grpc.NewServer(opt...)
	for _, service := range services {
		service.Register(server)
	}

	return &Server{
		server: server,
		port:   port,
	}, nil
}

// Run starts the server.
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return s.server.Serve(listener)
}

// Shutdown performs a graceful shutdown of the server where the given context can be used to enforce a timeout.
// The call blocks until this is done.
func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()

	case <-stopped:
		s.server.Stop()
		return nil
	}
}
