package main

import (
	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/server/rest"
	"fibonacci_service/pkg/server/rpc"
	"fibonacci_service/pkg/service"
	"log"
)

func main() {
	c := cache.New()
	svc := service.New(c)
	restServ := rest.New(svc)

	go func() {
		log.Fatal(restServ.StartRest(1323))
	}()

	rpcServ := rpc.New(svc)

	go func() {
		log.Fatal(rpcServ.StartRpc(9000))
	}()

	wait := make(chan chan bool)
	<-wait
}
