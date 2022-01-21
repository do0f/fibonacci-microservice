package main

import (
	"context"
	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/server/rest"
	"fibonacci_service/pkg/server/rpc"
	"fibonacci_service/pkg/service"
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
	restServ := rest.New(svc)

	restPort, err := strconv.Atoi(os.Getenv("REST_PORT"))
	if err != nil {
		log.Fatalf("invalid rest port number")
	}

	go func() {
		log.Fatal(restServ.StartRest(restPort))
	}()

	rpcServ := rpc.New(svc)

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("invalid grpc port number")
	}

	go func() {
		log.Fatal(rpcServ.StartRpc(grpcPort))
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Printf("attempting graceful shutdown")

	err = c.GracefulShutdown()
	if err != nil {
		log.Printf("failed to gracefully shutdown redis cache: %s", err.Error())
	}
	log.Printf("cache shutted down")

	rpcServ.GracefulShutdown()
	log.Printf("rpc server shutted down")

	err = restServ.GracefulShutdown(context.Background())
	if err != nil {
		log.Printf("failed to gracefully shutdown rest server: %s", err.Error())
	}
	log.Printf("rest server shutted down")
}
