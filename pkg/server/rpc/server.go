package rpc

import (
	"fibonacci_service/pkg/rpc"
	"fibonacci_service/pkg/server"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	s *grpc.Server
	rpc.UnimplementedFibonacciServer

	svc server.FibService
}

func New(service server.FibService) *Server {
	serv := new(Server)

	serv.svc = service

	serv.s = grpc.NewServer()
	rpc.RegisterFibonacciServer(serv.s, serv)
	return serv
}

func (serv *Server) StartRpc(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return serv.s.Serve(lis)
}
