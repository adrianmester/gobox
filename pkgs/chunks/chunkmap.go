package chunks

import (
	"crypto/md5"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type ChunkMap struct {
	Chunks map[string]ChunkID
	Mu     *sync.Mutex
}

func NewChunkMap() ChunkMap {
	return ChunkMap{
		Chunks: map[string]ChunkID{},
		Mu:     &sync.Mutex{},
	}
}

func (cm *ChunkMap) AddChunk(buf []byte, chunkID ChunkID) (ChunkID, bool) {
	checksum := fmt.Sprintf("%x", md5.Sum(buf))
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	// if the chunk is already in the map, return the existing ChunkID
	if existingChunkID, ok := cm.Chunks[checksum]; ok {
		return existingChunkID, true
	}
	// otherwise, add it to the map, and return the chunkID that was passed in the function call
	cm.Chunks[checksum] = chunkID
	return chunkID, false
}

const ChunkSize = 1024

type Chunk struct {
	ChunkID
	Data []byte
}

func (cm ChunkMap) GetFileChunks(baseDir string, fileID int64, path string) chan Chunk {
	result := make(chan Chunk)
	go func() {
		defer close(result)
		fp, err := os.Open(filepath.Join(baseDir, path))
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("failed to read file")
			return
		}
		var chunkNumber int64
		for {
			buffer := make([]byte, ChunkSize)
			bytesRead, err := fp.Read(buffer)
			if err != nil {
				if err == io.EOF {
					// to handle empty files, send one last chunk, without data, but with the fileId at the end
					chunk := Chunk{
						ChunkID: ChunkID{fileID, -1},
					}
					result <- chunk
					return
				}
				log.Error().Err(err).Str("path", path).Msg("failed to read file")
				return
			}

			chunkID := ChunkID{fileID, chunkNumber}
			chunkID, alreadyExists := cm.AddChunk(buffer, chunkID)

			chunk := Chunk{
				ChunkID: chunkID,
			}
			if !alreadyExists {
				chunk.Data = buffer[:bytesRead]
			}
			result <- chunk

			chunkNumber++
		}
	}()
	return result
}
