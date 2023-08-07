package parser

import (
	"fmt"
	"io"
)

func Parse(r io.Reader) ([]Chunk, error) {
	return ParseWithOptions(r, &ChunkParseOptions{}, true)
}

func ParseWithOptions(r io.Reader, options *ChunkParseOptions, unsafeByteToString bool) ([]Chunk, error) {
	var chunks []Chunk
	for {
		var chunk Chunk
		err := chunk.Parse(r, options, unsafeByteToString)
		if err == io.EOF {
			return chunks, nil
		}
		if err != nil {
			return chunks, fmt.Errorf("unable to parse chunk: %w", err)
		}
		chunks = append(chunks, chunk)
	}
}
