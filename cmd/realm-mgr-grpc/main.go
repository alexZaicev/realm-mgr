package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type exitCode int

const (
	exitCodeNoError = iota
	exitCodeConfigNotFound
	exitCodeAppInitError
	exitCodeGRPCServerRunError
	exitCodeGRPCServerShutdownError
)

func run() exitCode {
	ctx := context.Background()

	cfgFilePath, err := getConfigFile()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return exitCodeConfigNotFound
	}

	app, err := initialize(ctx, cfgFilePath)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return exitCodeAppInitError
	}

	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the GRPC server
	go func() {
		fmt.Println("INFO: Starting gRPC server...")

		if serverRunErr := app.grpcServer.Run(); serverRunErr != nil {
			fmt.Printf("ERROR: %s\n", serverRunErr)
			os.Exit(exitCodeGRPCServerRunError)
		}
	}()

	// No need for a select around this channel read, as it will block
	// until an interrupt signal is put on it
	<-endChan

	fmt.Println("INFO: Killing gRPC server....")

	if shutdownErr := app.grpcServer.Shutdown(ctx); shutdownErr != nil {
		fmt.Printf("ERROR: %s\n", shutdownErr)
		os.Exit(exitCodeGRPCServerShutdownError)
	}

	fmt.Println("INFO: gRPC server killed.")

	return exitCodeNoError
}

func main() {
	os.Exit(int(run()))
}
