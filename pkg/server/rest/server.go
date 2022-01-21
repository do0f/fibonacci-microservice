package rest

import (
	"context"
	"fibonacci_service/pkg/service"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

const (
	GetFibbonaciEndpoint = "/fibonacci"
)

type FibService interface {
	FibSequence(first int, last int) ([]service.FibNumber, error)
}

type Server struct {
	e   *echo.Echo
	svc FibService
}

func New(service FibService) *Server {
	serv := new(Server)

	serv.e = echo.New()
	serv.svc = service

	serv.e.GET(GetFibbonaciEndpoint, serv.GetFibonacciHandler)

	return serv
}

func (serv *Server) StartRest(port int) error {
	return serv.e.Start(fmt.Sprintf(":%d", port))
}

func (serv *Server) NewContext(h *http.Request, w http.ResponseWriter) echo.Context {
	return serv.e.NewContext(h, w)
}
func (serv *Server) GracefulShutdown(ctx context.Context) error {
	return serv.e.Shutdown(ctx)
}
