package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"io/fs"
	"os"
	"path/filepath"
)

func listDirectories(path string) ([]string, error) {
	result := []string{}
	wErr := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		if err != nil {
			return nil
		}
		result = append(result, path)
		return nil
	})
	if wErr != nil {
		return result, wErr
	}
	return result, nil
}

func watch(log zerolog.Logger, ctx context.Context, path string) (chan string, error) {
	result := make(chan string)

	watchedDirectories, err := listDirectories(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to list directories to watch")
		return result, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize file watcher")
		return result, err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				log.Debug().Str("event", event.Op.String()).Str("path", event.Name).Msg("received event")
				result <- event.Name
				lstat, err := os.Lstat(event.Name)
				if err != nil {
					continue
				}
				if !lstat.IsDir() {
					continue
				}
				err = watcher.Add(event.Name)
				if err != nil {
					log.Error().Err(err).Msg("couldn't add dir")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msg("event error")
			case <-ctx.Done():
				log.Info().Msg("shutting down the file watcher")
				close(result)
				watcher.Close()
				return
			}
		}
	}()

	for _, dir := range watchedDirectories {
		err = watcher.Add(dir)
		if err != nil {
			log.Error().Err(err).Str("path", dir).Msg("couldn't add path")
		}
	}
	return result, nil
}
