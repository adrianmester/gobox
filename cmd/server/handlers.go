package main

import (
	"context"
	"fmt"
	"github.com/adrianmester/gobox/proto"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func (g *goboxServer) SendFileInfo(_ context.Context, fileInfo *proto.SendFileInfoInput) (*proto.SendFileInfoResponse, error) {
	g.lock.Lock()
	g.Files[fileInfo.FileId] = File{
		Name:        fileInfo.FileName,
		IsDirectory: fileInfo.IsDirectory,
		ModTime:     time.Unix(fileInfo.ModTime, 0),
		Size:        fileInfo.Size,
	}
	g.lock.Unlock()

	fullPath := filepath.Join(g.dataDir, fileInfo.FileName)
	if fileInfo.IsDirectory {
		stat, err := os.Lstat(fullPath)
		if err == nil && !stat.IsDir() {
			// the path exists, but it's a file, we'll need to remove it first
			err = os.RemoveAll(fullPath)
			if err != nil {
				g.log.Error().Err(err).Str("path", fileInfo.FileName).Msg("failed to remove file")
			}
		}
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			g.log.Error().Err(err).Str("path", fileInfo.FileName).Msg("failed to create directory")
		}
		return &proto.SendFileInfoResponse{SendChunkIds: false}, nil
	} else {
		stat, err := os.Lstat(fullPath)
		if err == nil && stat.IsDir() {
			// the path exists, but it's a directory not a file, we'll need to remove it first
			err = os.RemoveAll(fullPath)
			if err != nil {
				g.log.Error().Err(err).Str("path", fileInfo.FileName).Msg("failed to remove directory")
			}
		}
	}

	return &proto.SendFileInfoResponse{SendChunkIds: g.DoesFileNeedUpdate(fileInfo.FileId)}, nil
}

func (g *goboxServer) SendFileChunks(stream proto.GoBox_SendFileChunksServer) error {
	var (
		chunkCount int64
		fileID     int64 = -1
		file       File
		fp         *os.File
	)
	defer func() {
		g.log.Debug().Str("path", file.Name).Int64("chunks", chunkCount).Msg("wrote file")
		err := os.Chtimes(filepath.Join(g.dataDir, file.Name), file.ModTime, file.ModTime)
		if err != nil {
			g.log.Error().Err(err).Str("path", file.Name).Msg("failed to update mtime")
		}
		_ = fp.Close()
	}()
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			g.log.Debug().
				Int64("chunks", chunkCount).
				Int64("fileID", fileID).
				Msg("received chunks")
			return nil
		}
		if err != nil {
			//TODO:
			g.log.Error().Err(err).Int64("fileID", fileID).Msg("not nil err")
			return nil
		}
		if fileID == -1 {
			// this is the first chunk, we need to do some initialisations
			fileID = chunk.ChunkId.FileId
			g.lock.Lock()
			file = g.Files[fileID]
			g.lock.Unlock()
			fp, err = os.Create(filepath.Join(g.dataDir, file.Name))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", file.Name, err)
			}
		}
		if len(chunk.Data) > 0 {
			_, err = fp.Write(chunk.Data)
			if err != nil {
				g.log.Error().Err(err).Msg("error writing file")
			}
			chunkCount += 1
		}
	}
}

func (g *goboxServer) InitialSyncComplete(context.Context, *proto.Null) (*proto.Null, error) {
	wantedPaths := map[string]bool{}
	g.lock.Lock()
	for _, file := range g.Files {
		wantedPaths[filepath.Join(g.dataDir, file.Name)] = true
	}
	g.lock.Unlock()
	err := filepath.Walk(g.dataDir, func(path string, info fs.FileInfo, err error) error {
		path = filepath.Clean(path)
		if _, ok := wantedPaths[path]; !ok {
			// this file wasn't one of the ones sent by the client
			g.log.Debug().Str("path", path).Msg("file not send by client, deleting")
			err := os.RemoveAll(path)
			if err != nil {
				g.log.Error().Err(err).Str("path", path).Msg("couldn't delete file")
			}
		}
		return nil
	})
	if err != nil {
		g.log.Error().Err(err).Msg("InitialSyncComplete walk")
	}
	g.log.Info().Msg("Initial Sync Complete")
	return &proto.Null{}, nil
}
func (g *goboxServer) DeleteFile(ctx context.Context, in *proto.DeleteFileInput) (*proto.Null, error) {
	err := os.RemoveAll(filepath.Join(g.dataDir, in.Path))
	return &proto.Null{}, err
}