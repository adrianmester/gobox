package chunks

import (
	"fmt"
	"os"
)

func GetChunk(path string, chunkSize, chunkNumber int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return []byte{}, fmt.Errorf("error getting chunk %d from %s: %w", chunkNumber, path, err)
	}
	var buffer = make([]byte, chunkSize)
	bytesRead, err := f.ReadAt(buffer, chunkNumber*chunkSize)
	if err != nil {
		return []byte{}, fmt.Errorf("error getting chunk %d from %s: %w", chunkNumber, path, err)
	}
	return buffer[:bytesRead], nil
}
