package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

type RegServiceFn func(s *grpc.Server)

type Server struct {
	grpcServer *grpc.Server

	regServiceFns []RegServiceFn
}

func NewServer(opt ...grpc.ServerOption) *Server {
	return &Server{
		grpcServer:    grpc.NewServer(opt...),
		regServiceFns: make([]RegServiceFn, 0),
	}
}

func (s *Server) applyServices() {
	for _, fn := range s.regServiceFns {
		fn(s.grpcServer)
	}
}

func (s *Server) ln(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func (s *Server) listenSigs() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

  <- ch

	s.GracefulStop()
}

// API
func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) AppendRegServiceFn(fn RegServiceFn) {
	s.regServiceFns = append(s.regServiceFns, fn)
}

func (s *Server) Serve(addr string) error {
	var (
		ln  net.Listener
		err error
	)

	if ln, err = s.ln(addr); err != nil {
		return fmt.Errorf("can not make listener, err: %s", err.Error())
	}

	go s.listenSigs()

	s.applyServices()

	err = s.grpcServer.Serve(ln)

	if err != nil {
		return err
	}

	return nil
}
