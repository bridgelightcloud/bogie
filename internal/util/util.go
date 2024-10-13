package util

import "github.com/bridgelightcloud/bogie/internal/models"

type ChunkifySliceParams[M models.Model] struct {
	Models    []M
	ChunkSize int
}

func ChunkifySlice[M models.Model](params ChunkifySliceParams[M]) [][]M {
	if len(params.Models) == 0 {
		return [][]M{}
	}

	if params.ChunkSize == 0 {
		params.ChunkSize = 25
	}

	chunks := make([][]M, 0)
	for i := 0; i < len(params.Models); i += params.ChunkSize {
		end := i + params.ChunkSize

		if end > len(params.Models) {
			end = len(params.Models)
		}

		chunks = append(chunks, params.Models[i:end])
	}

	return chunks
}
