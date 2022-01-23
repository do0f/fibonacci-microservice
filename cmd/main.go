package main

import (
	"context"
	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/server/rest"
	"fibonacci_service/pkg/server/rpc"
	"fibonacci_service/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"log"
)

func main() {
	cDb, err := strconv.Atoi(os.Getenv("CACHE_DB_NUM"))
	if err != nil {
		log.Fatalf("invalid cache db number")
	}
	c, err := cache.Connect(
		os.Getenv("CACHE_ADDRESS"),
		os.Getenv("CACHE_PASSWORD"),
		cDb,
	)
	if err != nil {
		log.Fatalf("failed to connect to redis cache: %s", err.Error())
	}

	svc := service.New(c)

	restPort, err := strconv.Atoi(os.Getenv("REST_PORT"))
	if err != nil {
		log.Fatalf("invalid rest port number")
	}
	restServ := rest.New(svc)
	go func() {
		if err := restServ.StartRest(restPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down rest server: %s", err.Error())
		}
	}()

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("invalid grpc port number")
	}
	rpcServ := rpc.New(svc)
	go func() {
		if err := rpcServ.StartRpc(grpcPort); err != nil {
			log.Fatalf("shutting down rpc server: %s", err.Error())
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Printf("attempting graceful shutdown")

	rpcServ.GracefulShutdown()
	log.Printf("rpc server shutted down")

	err = restServ.GracefulShutdown(context.Background())
	if err != nil {
		log.Printf("failed to gracefully shutdown rest server: %s", err.Error())
	}
	log.Printf("rest server shutted down")

	err = c.GracefulShutdown()
	if err != nil {
		log.Printf("failed to gracefully shutdown redis cache: %s", err.Error())
	}
	log.Printf("cache shutted down")
}
