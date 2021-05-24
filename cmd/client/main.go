package main

import (
	"context"
	"flag"
	"github.com/adrianmester/gobox/pkgs/logging"
	"github.com/adrianmester/gobox/proto"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	ctx, cancelMainContext := context.WithCancel(context.Background())
	defer cancelMainContext()
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

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := proto.NewGoBoxClient(conn)

	pm := newPathIDMap()
	wg := sync.WaitGroup{}

	// Initial sync
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

	// watch for file changes and keep the server updated
	for fileChanged := range fileChangesChan {
		fileID := pm.GetID(fileChanged)
		fInfo, err := os.Lstat(filepath.Join(dataDir, fileChanged))
		if err != nil {
			// file was removed on the client side, delete it on the server as well
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
}
