package main

import (
	"log"
	"path"
	"runtime"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/interceptors"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/config"
	"github.com/alexZaicev/realm-mgr/internal/drivers/grpcserver"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
)

const (
	ConfigLogLevel  = "log.level"
	ConfigGRPCPort  = "grpc.port"
	ConfigDBHost    = "database.host"
	ConfigDBPort    = "database.port"
	ConfigDBUser    = "database.user"
	ConfigDBPass    = "database.password"
	ConfigDBName    = "database.name"
	ConfigDBSSLMode = "database.ssl_mode"
)

type application struct {
	grpcServer *grpcserver.Server
}

func newApplication(grpcServer *grpcserver.Server) *application {
	return &application{
		grpcServer: grpcServer,
	}
}

func newConfigStore(cfgFilePath string) (config.Config, error) {
	cfg, err := config.NewKoanfConfig(".").
		WithTOML(cfgFilePath, false /*optional*/).
		TreatInt64AsInt(true).
		Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func newLoggerFromConfig(cfg config.Config) (logging.Logger, error) {
	logLevel, err := config.Get[string](cfg, ConfigLogLevel)
	if err != nil {
		return nil, err
	}
	return logging.NewZapJSONLogger(logLevel)
}

func newDBConnectionProviderFromConfig(cfg config.Config, options ...pgdb.DBOption) (*pgdb.ConnProvider, error) {
	dbHost, err := config.Get[string](cfg, ConfigDBHost)
	if err != nil {
		return nil, err
	}
	dbPort, err := config.Get[string](cfg, ConfigDBPort)
	if err != nil {
		return nil, err
	}
	dbUser, err := config.Get[string](cfg, ConfigDBUser)
	if err != nil {
		return nil, err
	}
	dbPassword, err := config.Get[string](cfg, ConfigDBPass)
	if err != nil {
		return nil, err
	}
	dbName, err := config.Get[string](cfg, ConfigDBName)
	if err != nil {
		return nil, err
	}
	dbSSLMode, err := config.Get[string](cfg, ConfigDBSSLMode)
	if err != nil {
		return nil, err
	}

	return pgdb.NewConnProvider(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode, options...)
}

func newDBOptions() ([]pgdb.DBOption, error) {
	return nil, nil
}

func newGRPCServerFromConfig(
	cfg config.Config,
	services []grpcserver.Service,
	opt ...grpc.ServerOption,
) (*grpcserver.Server, error) {
	port, err := config.Get[int](cfg, ConfigGRPCPort)
	if err != nil {
		return nil, err
	}

	return grpcserver.NewServer(uint16(port), services, opt...)
}

func newGRPCServices(realmMgrAPI *realmmgrgrpc.RealmManagerAPI) ([]grpcserver.Service, error) {
	return []grpcserver.Service{
		realmMgrAPI,
	}, nil
}

func newGRPCServerOptions(logger logging.Logger) ([]grpc.ServerOption, error) {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptors.LoggerUnaryServerInterceptor(logger),
			interceptors.ValidateUnaryServerInterceptor(logger),
		)),
	}, nil
}

func getConfigFile() (string, error) {
	_, currentFileName, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("error fetching config file Name")
		return "", realmmgr_errors.NewInternalError("Failed to fetch config file name", nil)
	}
	currentFilePath := path.Dir(currentFileName)
	cfgFile := currentFilePath + "/../../config.toml"
	return cfgFile, nil
}
