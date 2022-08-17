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
	_grpcServer *grpc.Server

	regServiceFns  []RegServiceFn
	interceptedSig os.Signal
}

func NewServer(opt ...grpc.ServerOption) *Server {
	return &Server{
		_grpcServer:    grpc.NewServer(opt...),
		regServiceFns:  make([]RegServiceFn, 0),
		interceptedSig: nil,
	}
}

func (s *Server) _ln(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func (s *Server) _appendRegServiceFn(fn RegServiceFn) {
	s.regServiceFns = append(s.regServiceFns, fn)
}

func (s *Server) _listenSigs() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	s.interceptedSig = <-ch

	s._stop()
}

func (s *Server) _stop() {
	s._grpcServer.GracefulStop()
}

func (s *Server) _serve(addr string) error {
	var (
		ln  net.Listener
		err error
	)

	if ln, err = s._ln(addr); err != nil {
		return fmt.Errorf("can not make listener, err: %s", err.Error())
	}

	// listen sigs
	go s._listenSigs()

	err = s._grpcServer.Serve(ln)

	if err != nil {
		return err
	}

	if s.interceptedSig != nil {
		return fmt.Errorf("intercepted sig: %s", s.interceptedSig.String())
	}

	return nil
}

// API
func (s *Server) GracefulStop() {
	s._stop()
}

func (s *Server) AppendRegServiceFn(fn RegServiceFn) {
	s._appendRegServiceFn(fn)
}

func (s *Server) Serve(addr string) error {
	return s._serve(addr)
}
