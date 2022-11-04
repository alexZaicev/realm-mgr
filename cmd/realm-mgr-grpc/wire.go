//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	adaptercommon "github.com/alexZaicev/realm-mgr/internal/adapters/common"
	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc"
	"github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
	"github.com/alexZaicev/realm-mgr/internal/usecases/realms"
)

func initialize(ctx context.Context, cfgFilePath string) (*application, error) {
	panic(wire.Build(
		// UUID generator
		uuidgenerator.NewGoogleUUIDGenerator,
		wire.Bind(new(uuidgenerator.Generator), new(uuidgenerator.GoogleUUIDGenerator)),
		// Clock
		clock.NewStdLibClock,
		wire.Bind(new(clock.Clock), new(clock.StdLibClock)),
		// Configuration
		newConfigStore,
		// Logger
		newLoggerFromConfig,
		// Repositories
		newDBOptions,
		newDBConnectionProviderFromConfig,
		wire.Bind(new(postgres.ReadConnProvider), new(*pgdb.ConnProvider)),
		wire.Bind(new(postgres.WriteConnProvider), new(*pgdb.ConnProvider)),
		postgres.NewDataStoreLifecycleManager,
		adaptercommon.NewPgDataStoreManager,
		wire.Bind(new(adaptercommon.PgDatastoreLifeCycleManager), new(*postgres.DataStoreLifecycleManager)),
		wire.Bind(new(adaptercommon.DataStoreManager), new(*adaptercommon.PgDataStoreManager)),
		// UseCases
		realms.NewGetRealm,
		realms.NewCreateRealm,
		realms.NewReleaseRealm,
		// UseCase executors
		wire.Bind(new(adaptercommon.RealmGetter), new(*realms.GetRealm)),
		wire.Bind(new(adaptercommon.RealmCreator), new(*realms.CreateRealm)),
		wire.Bind(new(adaptercommon.RealmReleaser), new(*realms.ReleaseRealm)),
		adaptercommon.NewRealmUseCaseExecutor,
		// gRPC server
		wire.Bind(new(realmmgrgrpc.RealmOps), new(*adaptercommon.RealmUseCaseExecutor)),
		realmmgrgrpc.NewRealmManagerAPI,
		newGRPCServices,
		newGRPCServerOptions,
		newGRPCServerFromConfig,
		// main
		newApplication,
	))
}
