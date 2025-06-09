package main

import (
	"api-standardisation/grpc_bro"
	notesapi_v1 "api-standardisation/openapi3"
	"api-standardisation/restapi"
	"api-standardisation/store"
	"api-standardisation/tsp-output/notesapi"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dataStore := store.NewInMemoryStore()

	errChan := make(chan error, 2)

	go func() {
		errChan <- startHTTPServer(dataStore)
	}()
	go func() {
		errChan <- startGRPCServer(dataStore)
	}()

	// Wait for interruption signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for server error or interruption
	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	log.Println("Shutting down servers...")
}
func startHTTPServer(store store.Store) error {
	e := echo.New()
	e.Use((middleware.Logger()))
	e.Use(middleware.Recover())

	httpHandler := restapi.NewHTTPServer(store)
	notesapi_v1.RegisterHandlers(e, httpHandler)
	// Start server
	srv := &http.Server{
		Addr:    ":8093",
		Handler: e,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Println("HTTP server listening on :8093")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
func startGRPCServer(dataStore store.Store) error {
	// Create gRPC server
	grpcServer := grpc.NewServer()
	serviceImpl := grpc_bro.NewGRPCServer(dataStore)
	// Register services
	notesapi.RegisterAuthServer(grpcServer, serviceImpl)
	notesapi.RegisterNotesServer(grpcServer, serviceImpl)

	// Start listening
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}

	log.Println("gRPC server listening on :50052")
	return grpcServer.Serve(lis)
}
