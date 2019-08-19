package internal_api

import (
	"log"
	"net"

	"github.com/Quard/poindexter/internal/storage"
	grpc "google.golang.org/grpc"
)

type Opts struct {
	bind string
}

type internalAPIServer struct {
	opts       Opts
	listener   net.Listener
	grpcServer *grpc.Server
	storage    storage.Storage
}

func NewInternalAPIServer(bind string, stor storage.Storage) internalAPIServer {
	var err error

	srv := internalAPIServer{
		opts:    Opts{bind: bind},
		storage: stor,
	}
	srv.listener, err = net.Listen("tcp", srv.opts.bind)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	srv.grpcServer = grpc.NewServer(opts...)
	RegisterInternalAPIServer(srv.grpcServer, &srv)

	return srv
}

func (srv internalAPIServer) Run() {
	log.Printf("INTERNAL gRPC server listen on: %s", srv.opts.bind)
	srv.grpcServer.Serve(srv.listener)
}
