package main

import (
	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/server"
	"fibonacci_service/pkg/service"
	"log"
)

func main() {
	c := cache.New()
	svc := service.New(c)
	serv := server.New(svc)

	if err := serv.Start(1323); err != nil {
		log.Fatal(err.Error())
	}
}
