package gollama

import (
	"context"
	"math"
)

// Embedding generates a vector embedding for a given string of text using the
// currently set model. The model must support the "embeddings" capability.
//
// The function will return an error if the model does not support the
// "embeddings" capability. The function will also return an error if the
// request fails.
//
// The function returns a slice of floats, representing the vector
// embedding of the input text.
func (c *Gollama) Embedding(ctx context.Context, prompt string) ([]float64, error) {
	req := embeddingsRequest{
		Model:  c.ModelName,
		Prompt: prompt,
	}

	var resp embeddingsResponse
	err := c.apiPost(ctx, "/api/embeddings", &resp, req)
	if err != nil {
		return nil, err
	}

	return resp.Embedding, nil
}

func CosenoSimilarity(vector1, vector2 []float64) float64 {
	if len(vector1) != len(vector2) {
		return 0.0
	}

	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for i := 0; i < len(vector1); i++ {
		dotProduct += vector1[i] * vector2[i]
		norm1 += vector1[i] * vector1[i]
		norm2 += vector2[i] * vector2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}
