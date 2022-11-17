package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	flag "github.com/spf13/pflag"

	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/interceptors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/config"
	"github.com/alexZaicev/realm-mgr/internal/drivers/grpcserver"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
)

const (
	FlagConfig      = "config"
	FlagConfigShort = "c"
)

const (
	configLogLevel  = "logging.level"
	configGRPCPort  = "grpc.port"
	configDBHost    = "database.host"
	configDBPort    = "database.port"
	configDBUser    = "database.user"
	configDBPass    = "database.password"
	configDBName    = "database.name"
	configDBSSLMode = "database.ssl_mode"
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
		WithYAML(cfgFilePath, false /*optional*/).
		TreatInt64AsInt(true).
		Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func newLoggerFromConfig(cfg config.Config) (logging.Logger, error) {
	logLevel, err := config.Get[string](cfg, configLogLevel)
	if err != nil {
		return nil, err
	}
	return logging.NewZapJSONLogger(logLevel)
}

func newDBConnectionProviderFromConfig(cfg config.Config, options ...pgdb.DBOption) (*pgdb.ConnProvider, error) {
	dbHost, err := config.Get[string](cfg, configDBHost)
	if err != nil {
		return nil, err
	}
	dbPort, err := config.Get[int](cfg, configDBPort)
	if err != nil {
		return nil, err
	}
	dbUser, err := config.Get[string](cfg, configDBUser)
	if err != nil {
		return nil, err
	}
	dbPassword, err := config.Get[string](cfg, configDBPass)
	if err != nil {
		return nil, err
	}
	dbName, err := config.Get[string](cfg, configDBName)
	if err != nil {
		return nil, err
	}
	dbSSLMode, err := config.Get[string](cfg, configDBSSLMode)
	if err != nil {
		return nil, err
	}

	return pgdb.NewConnProvider(dbHost, fmt.Sprintf("%d", dbPort), dbUser, dbPassword, dbName, dbSSLMode, options...)
}

func newDBOptions() ([]pgdb.DBOption, error) {
	return nil, nil
}

func newGRPCServerFromConfig(
	cfg config.Config,
	services []grpcserver.Service,
	opt ...grpc.ServerOption,
) (*grpcserver.Server, error) {
	port, err := config.Get[int](cfg, configGRPCPort)
	if err != nil {
		return nil, err
	}

	return grpcserver.NewServer(uint16(port), services, opt...)
}

func newGRPCServices(realmMgrAPI *realmmgrgrpc.RealmManagerAPI) []grpcserver.Service {
	return []grpcserver.Service{
		realmMgrAPI,
	}
}

func newGRPCServerOptions(logger logging.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptors.LoggerUnaryServerInterceptor(logger),
			interceptors.ValidateUnaryServerInterceptor(logger),
		)),
	}
}

func getConfigFile() (string, error) {
	args := os.Args[1:]

	f := flag.NewFlagSet(FlagConfig, flag.ContinueOnError)
	f.StringP(FlagConfig, FlagConfigShort, "", "Service YAML configuration file")
	if err := f.Parse(args); err != nil {
		return "", fmt.Errorf("error parsing flags: %w", err)
	}

	cfgFile, err := f.GetString(FlagConfig)
	if err != nil {
		return "", fmt.Errorf("error reading flags: %w", err)
	}
	return cfgFile, nil
}
