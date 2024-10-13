package util

import (
	"math"
)

func ChunkifySlice[T any](items []T, chunkSize int) [][]T {
	if len(items) <= chunkSize {
		return [][]T{items}
	}

	chunkCount := int(math.Ceil(float64(len(items)) / float64(chunkSize)))
	chunks := make([][]T, chunkCount)

	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize

		if end > len(items) {
			end = len(items)
		}

		chunks = append(chunks, items[i:end])
	}

	return chunks
}
