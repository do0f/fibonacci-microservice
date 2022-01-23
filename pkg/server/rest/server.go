package rest

import (
	"context"
	"fibonacci_service/pkg/server"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	GetFibbonaciEndpoint = "/fibonacci"
)

type Server struct {
	e   *echo.Echo
	svc server.FibService
}

type Context struct {
	echo.Context
}

func New(service server.FibService) *Server {
	serv := new(Server)

	serv.e = echo.New()
	serv.svc = service

	serv.e.GET(GetFibbonaciEndpoint, serv.GetFibonacciHandler)

	return serv
}

func (serv *Server) StartRest(port int) error {
	serv.e.Logger.SetLevel(log.INFO)
	return serv.e.Start(fmt.Sprintf(":%d", port))
}

func (serv *Server) NewContext(h *http.Request, w http.ResponseWriter) Context {
	return Context{serv.e.NewContext(h, w)}
}
func (serv *Server) GracefulShutdown(ctx context.Context) error {
	return serv.e.Shutdown(ctx)
}
