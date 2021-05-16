package main

import (
	"context"
	"fmt"
	"github.com/adrianmester/gobox/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	lis, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		log.Fatal().Err(err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterGoBoxServer(grpcServer, NewGoBoxServer())
	grpcServer.Serve(lis)
}

type goboxServer struct {
	proto.UnimplementedGoBoxServer

	Files map[int64]string
}
func NewGoBoxServer() *goboxServer {
	return &goboxServer{
		Files: map[int64]string{},
	}
}

func (g *goboxServer) GetLastUpdateTime(_ context.Context, _ *proto.Null) (*proto.GetLastUpdateTimeResult, error) {
	result := proto.GetLastUpdateTimeResult{
		Timestamp: time.Now().Unix(),
	}
	return &result, nil
}

func (g *goboxServer) SendFileInfo(_ context.Context, fileInfo *proto.SendFileInfoInput) (*proto.SendFileInfoResponse, error) {
	g.Files[fileInfo.FileId] = fileInfo.FileName
	fmt.Println(g.Files)
	return &proto.SendFileInfoResponse{SendChunkIds: true}, nil
}

func (g *goboxServer) SendFileChunks(stream proto.GoBox_SendFileChunksServer) error {
	chunkCount := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Debug().Int("chunks", chunkCount).Msg("received chunks")
			return nil
		}
		if err != nil {
			log.Panic().Err(err)
		}
		chunkCount += 1
	}
}
