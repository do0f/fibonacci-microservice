package server

import (
	"fibonacci_service/pkg/service"
	"fmt"

	echo "github.com/labstack/echo/v4"
)

const (
	GetFibbonaciEndpoint = "/fibonacci"
)

type IFibService interface {
	FibSequence(first int, last int) ([]service.FibNumber, error)
}

type Server struct {
	e   *echo.Echo
	svc IFibService
}

func New(service IFibService) *Server {
	serv := &Server{
		e:   echo.New(),
		svc: service,
	}

	serv.e.GET(GetFibbonaciEndpoint, serv.getFibonacci)

	return serv
}

func (serv *Server) Start(port int) error {
	return serv.e.Start(fmt.Sprintf(":%d", port))
}
