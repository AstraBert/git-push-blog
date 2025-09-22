package embeddings

import (
	"context"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/param"
)

func EmbedText(text string) ([]float64, error) {
	client := openai.NewClient()
	embd, err := client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Input:      openai.EmbeddingNewParamsInputUnion{OfString: param.Opt[string]{Value: text}},
		Model:      openai.EmbeddingModelTextEmbedding3Small,
		Dimensions: param.Opt[int64]{Value: 768},
	})
	if err != nil {
		return nil, err
	}
	embds := [][]float64{}
	for _, item := range embd.Data {
		embds = append(embds, item.Embedding)
	}
	return embds[0], nil
}
