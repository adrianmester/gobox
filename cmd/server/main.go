package main

import (
	"flag"
	"github.com/adrianmester/gobox/pkgs/logging"
	"github.com/adrianmester/gobox/proto"
	"google.golang.org/grpc"
	"net"
)

func main() {
	var (
		listenAddress string
		dataDir       string
		help          bool
		debug         bool
	)
	flag.StringVar(&listenAddress, "listen", "localhost:5555", "address to listen on (<host>:<port>)")
	flag.StringVar(&dataDir, "datadir", "./datadir/server", "path to directory to store files")
	flag.BoolVar(&help, "help", false, "show usage information")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	log := logging.GetLogger("server", debug)

	lis, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterGoBoxServer(grpcServer, NewGoBoxServer(&log, dataDir))
	log.Info().Str("address", listenAddress).Str("datadir", dataDir).Msg("Starting server")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Error().Err(err).Msg("Server error")
	}
}
