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
	"path/filepath"
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

type File struct {
	Name string
	IsDirectory bool
	ModTime time.Time
	Size int64
}

type goboxServer struct {
	proto.UnimplementedGoBoxServer

	Files map[int64]File
}
func NewGoBoxServer() *goboxServer {
	return &goboxServer{
		Files: map[int64]File{},
	}
}

func (g *goboxServer) GetLastUpdateTime(_ context.Context, _ *proto.Null) (*proto.GetLastUpdateTimeResult, error) {
	result := proto.GetLastUpdateTimeResult{
		//TODO:
		Timestamp: time.Now().Unix(),
	}
	return &result, nil
}

const DataDir="./datadir/server"

func (g *goboxServer) FileNeedsUpdate(fileID int64) bool{
	file := g.Files[fileID]
	fullPath := filepath.Join(DataDir, file.Name)
	fInfo, err := os.Stat(fullPath)
	if err != nil {
		// file doesn't exist
		log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Msg("file doesn't exist")
		return true
	}
	if fInfo.Size() != file.Size {
		log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Int64("expected size", file.Size).
			Int64("actual size", fInfo.Size()).
			Msg("file size doesn't match")
		return true
	}
	if fInfo.ModTime() != file.ModTime {
		log.Debug().Str("path", file.Name).
			Int64("fileID", fileID).
			Time("expected mtime", file.ModTime).
			Time("actual mtime", fInfo.ModTime()).
			Msg("file mtime doesn't match")
		return true
	}
	return false
}

func (g *goboxServer) SendFileInfo(_ context.Context, fileInfo *proto.SendFileInfoInput) (*proto.SendFileInfoResponse, error) {
	g.Files[fileInfo.FileId] = File{
		Name: fileInfo.FileName,
		IsDirectory: fileInfo.IsDirectory,
		ModTime: time.Unix(fileInfo.ModTime, 0),
		Size: fileInfo.Size,
	}

	if fileInfo.IsDirectory {
		log.Debug().Str("path", fileInfo.FileName).Msg("creating directory")
		err := os.MkdirAll(filepath.Join(DataDir, fileInfo.FileName), 0755)
		if err != nil {
			log.Error().Err(err).Str("path", fileInfo.FileName).Msg("failed to create directory")
		}
		return &proto.SendFileInfoResponse{SendChunkIds: false}, nil
	}

	return &proto.SendFileInfoResponse{SendChunkIds: g.FileNeedsUpdate(fileInfo.FileId)}, nil
}

func (g *goboxServer) SendFileChunks(stream proto.GoBox_SendFileChunksServer) error {
	var (
		chunkCount int64
		fileID int64 = -1
		file File
		fp *os.File
	)
	defer func(){
		log.Debug().Str("path", file.Name).Int64("chunks", chunkCount).Msg("wrote file")
		err := os.Chtimes(filepath.Join(DataDir, file.Name), file.ModTime, file.ModTime)
		if err != nil {
			log.Error().Err(err).Str("path", file.Name).Msg("failed to update mtime")
		}
		_ = fp.Close()
	}()
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			log.Debug().
				Int64("chunks", chunkCount).
				Int64("fileID", fileID).
				Msg("received chunks")
			return nil
		}
		if err != nil {
			//TODO:
			log.Error().Err(err).Int64("fileID", fileID).Msg("not nill err")
			return nil
		}
		if fileID == -1 {
			// this is the first chunk, we need to do some initialisations
			fileID = chunk.ChunkId.FileId
			file = g.Files[fileID]
			fp, err = os.Create(filepath.Join(DataDir, file.Name))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", file.Name, err)
			}
		}
		if len(chunk.Data) == 0 {
			log.Panic().Msg("missing chunk data, not implemented yet")
		}
		_, err = fp.Write(chunk.Data)
		if err != nil {
			log.Error().Err(err).Msg("error writing file")
		}
		chunkCount += 1
	}
}
