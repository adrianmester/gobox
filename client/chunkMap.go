package main

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

type ChunkID struct {
	FileID int64
	ChunkNumber int64
}

type ChunkMap struct {
	Chunks map[string]ChunkID
}

const ChunkLength = 1024 * 1024

type Chunk struct {
	ChunkID
	Data []byte
}

func (c ChunkMap) GetFileChunks(baseDir string, fileID int64, path string) chan Chunk {
	result := make(chan Chunk)
	go func() {
		defer close(result)
		fp, err := os.Open(filepath.Join(baseDir, path))
		if err != nil {
			//TODO:
			log.Panic().Err(err)
		}
		buffer := make([]byte, ChunkLength)
		var chunkNumber int64
		for {
			bytesRead, err := fp.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Error().Err(err).Str("path", path).Msg("failed to read file")
				return
			}
			buf := buffer[:bytesRead]
			//checksum := md5.Sum(buf)
			result <- Chunk{
				ChunkID{fileID, chunkNumber},
				buf,
			}

			chunkNumber++
		}
	}()
	return result
}
