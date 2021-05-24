package main

import (
	"context"
	"github.com/adrianmester/gobox/pkgs/chunks"
	"github.com/adrianmester/gobox/proto"
	"github.com/rs/zerolog"
	"io"
	"io/fs"
	"path/filepath"
	"sync"
)

type FileInfo struct {
	Path   string
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

func scanDirectory(log zerolog.Logger, pm *pathIDMap, dir string) chan FileInfo {
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

func sendChunksForFile(log zerolog.Logger, baseDir string, fInfo FileInfo, client *proto.GoBoxClient) {
	if fInfo.IsDir() {
		return
	}
	log.Debug().Str("path", fInfo.Path).Msg("sending chunks")
	cl, err := (*client).SendFileChunks(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("create send file chunks client")
		return
	}
	cm := chunks.NewChunkMap()
	for chunk := range cm.GetFileChunks(baseDir, fInfo.PathID, fInfo.Path) {
		err = cl.Send(&proto.SendFileChunksInput{
			ChunkId: &proto.ChunkID{
				ChunkNumber: chunk.ChunkNumber,
				FileId:      chunk.FileID,
			},
			Data: chunk.Data,
		})
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to send chunk")
		}
	}
	err = cl.CloseSend()
	if err != nil {
		log.Error().Err(err).Msg("close send")
	}
}

func updatePath(ctx context.Context, log zerolog.Logger, client *proto.GoBoxClient, wg *sync.WaitGroup, dataDir string, fInfo FileInfo) {
	log.Debug().Str("path", fInfo.Path).Int64("fileID", fInfo.PathID).Msg("sending path info")
	response, err := (*client).SendFileInfo(ctx, &proto.SendFileInfoInput{
		FileId:      fInfo.PathID,
		FileName:    fInfo.Path,
		IsDirectory: fInfo.IsDir(),
		Size:        fInfo.Size(),
		ModTime:     fInfo.ModTime().Unix(),
	})
	if err != nil {
		log.Error().Err(err).Msg("SendFileInfo")
		return
	}
	if response.SendChunkIds && !fInfo.IsDir() {
		wg.Add(1)
		go func(fInfo FileInfo) {
			defer wg.Done()
			sendChunksForFile(log, dataDir, fInfo, client)
		}(fInfo)
	}
}

