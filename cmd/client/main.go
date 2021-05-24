package main

import (
	"context"
	"flag"
	"github.com/adrianmester/gobox/pkgs/chunks"
	"github.com/adrianmester/gobox/pkgs/logging"
	"github.com/adrianmester/gobox/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"io"
	"io/fs"
	"os"
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

func main() {
	ctx, cancelMainContext := context.WithCancel(context.Background())
	var (
		serverAddress string
		dataDir       string
		help          bool
		debug         bool
	)
	flag.StringVar(&serverAddress, "server", "localhost:5555", "server address to connect to (<host>:<port>)")
	flag.StringVar(&dataDir, "datadir", "./datadir/client", "path to sync to remote server")
	flag.BoolVar(&help, "help", false, "show usage information")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	log := logging.GetLogger("client", debug)

	pm := newPathIDMap()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := proto.NewGoBoxClient(conn)
	wg := sync.WaitGroup{}

	for fInfo := range scanDirectory(log, pm, dataDir) {
		updatePath(ctx, log, &client, &wg, dataDir, fInfo)
	}
	wg.Wait()

	_, err = client.InitialSyncComplete(ctx, &proto.Null{})
	if err != nil {
		log.Error().Err(err).Msg("InitialSyncComplete")
	}
	log.Info().Msg("Initial Sync Complete")
	fileChangesChan, err := watch(log, ctx, dataDir)
	if err != nil {
		log.Error().Err(err).Msg("file watcher error")
	}
	for fileChanged := range fileChangesChan {
		fileID := pm.GetID(fileChanged)
		fInfo, err := os.Lstat(filepath.Join(dataDir, fileChanged))
		if err != nil {
			log.Debug().Str("path", fileChanged).Msg("deleting file")
			_, err := client.DeleteFile(ctx, &proto.DeleteFileInput{
				Path: fileChanged,
			})
			if err != nil {
				log.Error().Err(err).Str("path", fileChanged).Msg("couldn't delete remote file")
			}
			continue
		}
		updatePath(ctx, log, &client, &wg, dataDir, NewFileInfo(fileID, fileChanged, fInfo))
	}
	cancelMainContext()
}
