package main

import (
	"context"
	"github.com/adrianmester/gobox/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileInfo struct {
	Path string
	PathID int64
	fs.FileInfo
}

func NewFileInfo(fileID int64, path string, info fs.FileInfo) FileInfo {
	return FileInfo{
		path,
		fileID,
		info,
	}
}

func scanDirectory(pm *pathIDMap, dir string) chan FileInfo {
	dir = filepath.Clean(dir)
	log.Info().Str("path", dir).Msgf("scanning directory")
	result := make(chan FileInfo)
	go func() {
		defer close(result)
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Error().Err(err).Str("path", path).Msg("i/o error")
				return nil
			}
			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				log.Error().Err(err).Msg("get relative path")
				return nil
			}
			fileID := pm.GetID(relativePath)
			result <- NewFileInfo(fileID, relativePath, info)
			return nil
		})
		if err != nil {
			log.Error().Err(err)
		}
	}()
	return result
}

func sendChunksForFile(baseDir string, fInfo FileInfo, client *proto.GoBoxClient) {
	log.Debug().Str("path", fInfo.Path).Msg("sending chunks")
	cl, err := (*client).SendFileChunks(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("create send file chunks client")
		return
	}
	cm := ChunkMap{}
	if fInfo.IsDir() {
		return
	}
	for chunk := range cm.GetFileChunks(baseDir, fInfo.PathID, fInfo.Path) {
		err = cl.Send(&proto.SendFileChunksInput{
			ChunkId: &proto.ChunkID{
				ChunkNumber: chunk.ChunkNumber,
				FileId: chunk.FileID,
			},
			Data: chunk.Data,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to send chunk")
		}
	}
	_, err = cl.CloseAndRecv()
	if err != nil && err != io.EOF {
		log.Error().Err(err).Msg("close SendChunkIds")
		return
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	pm := NewPathIDMap()


	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost:5555", opts...)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := proto.NewGoBoxClient(conn)
	wg := sync.WaitGroup{}

	dataDir := "./datadir/client"
	for fInfo := range scanDirectory(pm, dataDir) {
		log.Debug().Str("path", fInfo.Path).Msg("sending file info")
		response, err := client.SendFileInfo(context.Background(), &proto.SendFileInfoInput{
			FileId: fInfo.PathID,
			FileName: fInfo.Path,
			IsDirectory: fInfo.IsDir(),
			Size: fInfo.Size(),
			ModTime: fInfo.ModTime().Unix(),
		})
		if err != nil {
			log.Error().Err(err).Msg("SendFileInfo")
			continue
		}
		if response.SendChunkIds {
			wg.Add(1)
			go func(fInfo FileInfo) {
				defer wg.Done()
				time.Sleep(time.Second)
				sendChunksForFile(dataDir, fInfo, &client)
			}(fInfo)
		}
	}
	wg.Wait()
	/*
	ts, err := client.GetLastUpdateTime(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Println(ts.Timestamp)

	client.SendFileInfo(context.Background(), &proto.SendFileInfoInput{
		FileName: "foo",
		FileId: 123,
	})
	*/
}
