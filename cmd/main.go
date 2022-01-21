package main

import (
	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/server/rest"
	"fibonacci_service/pkg/server/rpc"
	"fibonacci_service/pkg/service"
	"os"
	"strconv"

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

	wait := make(chan chan bool)
	<-wait
}
