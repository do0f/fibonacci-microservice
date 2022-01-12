package server

import (
	"fibonacci_service/pkg/service"
	"fmt"

	echo "github.com/labstack/echo/v4"
)

const (
	GetFibbonaciEndpoint = "/fibonacci"
)

type FibService interface {
	FibSequence(first int, last int) ([]service.FibNumber, error)
}

type Server struct {
	*echo.Echo
	svc FibService
}

func New(service FibService) *Server {
	serv := new(Server)
	serv.Echo = echo.New()
	serv.svc = service

	serv.GET(GetFibbonaciEndpoint, serv.GetFibonacci)

	return serv
}

func (serv *Server) Start(port int) error {
	return serv.Echo.Start(fmt.Sprintf(":%d", port))
}
