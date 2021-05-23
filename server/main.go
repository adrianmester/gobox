package main

import (
	"context"
	"flag"
	"github.com/adrianmester/gobox/pkgs/logging"
	"github.com/adrianmester/gobox/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
	"time"
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
	log := logging.GetLogger("client", debug)

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

type File struct {
	Name        string
	IsDirectory bool
	ModTime     time.Time
	Size        int64
}

type goboxServer struct {
	proto.UnimplementedGoBoxServer

	log     *zerolog.Logger
	dataDir string
	Files   map[int64]File
}

func NewGoBoxServer(log *zerolog.Logger, dataDir string) *goboxServer {
	return &goboxServer{
		Files:   map[int64]File{},
		dataDir: dataDir,
		log:     log,
	}
}

func (g *goboxServer) GetLastUpdateTime(_ context.Context, _ *proto.Null) (*proto.GetLastUpdateTimeResult, error) {
	result := proto.GetLastUpdateTimeResult{
		//TODO:
		Timestamp: time.Now().Unix(),
	}
	return &result, nil
}

func (g *goboxServer) DoesFileNeedUpdate(fileID int64) bool {
	file := g.Files[fileID]
	fullPath := filepath.Join(g.dataDir, file.Name)
	fInfo, err := os.Stat(fullPath)
	if err != nil {
		// file doesn't exist
		g.log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Msg("file doesn't exist")
		return true
	}
	if fInfo.Size() != file.Size {
		g.log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Int64("expected size", file.Size).
			Int64("actual size", fInfo.Size()).
			Msg("file size doesn't match")
		return true
	}
	if fInfo.ModTime() != file.ModTime {
		g.log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Time("expected mtime", file.ModTime).
			Time("actual mtime", fInfo.ModTime()).
			Msg("file mtime doesn't match")
		return true
	}
	return false
}

