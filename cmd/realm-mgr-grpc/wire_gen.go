// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"

	"github.com/alexZaicev/realm-mgr/internal/adapters/common"
	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc"
	"github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
	"github.com/alexZaicev/realm-mgr/internal/usecases/realms"
)

// Injectors from wire.go:

func initialize(ctx context.Context, cfgFilePath string) (*application, error) {
	config, err := newConfigStore(cfgFilePath)
	if err != nil {
		return nil, err
	}
	logger, err := newLoggerFromConfig(config)
	if err != nil {
		return nil, err
	}
	healthCheckService := realmmgrgrpc.NewHealthChecker(logger)
	googleUUIDGenerator := uuidgenerator.NewGoogleUUIDGenerator()
	stdLibClock := clock.NewStdLibClock()
	v, err := newDBOptions()
	if err != nil {
		return nil, err
	}
	connProvider, err := newDBConnectionProviderFromConfig(config, v...)
	if err != nil {
		return nil, err
	}
	dataStoreLifecycleManager, err := postgres.NewDataStoreLifecycleManager(connProvider, connProvider, googleUUIDGenerator)
	if err != nil {
		return nil, err
	}
	pgDataStoreManager, err := common.NewPgDataStoreManager(dataStoreLifecycleManager)
	if err != nil {
		return nil, err
	}
	getRealm := realms.NewGetRealm()
	createRealm := realms.NewCreateRealm()
	releaseRealm := realms.NewReleaseRealm()
	updateRealm := realms.NewUpdateRealm()
	realmUseCaseExecutor, err := common.NewRealmUseCaseExecutor(googleUUIDGenerator, stdLibClock, pgDataStoreManager, getRealm, createRealm, releaseRealm, updateRealm)
	if err != nil {
		return nil, err
	}
	realmManagerAPI, err := realmmgrgrpc.NewRealmManagerAPI(logger, realmUseCaseExecutor)
	if err != nil {
		return nil, err
	}
	v2 := newGRPCServices(healthCheckService, realmManagerAPI)
	v3 := newGRPCServerOptions(logger)
	server, err := newGRPCServerFromConfig(config, v2, v3...)
	if err != nil {
		return nil, err
	}
	mainApplication := newApplication(server)
	return mainApplication, nil
}
