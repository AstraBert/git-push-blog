package embeddings

import (
	"os"
	"testing"
)

func TestEmbedText(t *testing.T) {
	_, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		t.Skip("Skipping test because OPENAI_API_KEY is not in the environment")
	}
	embd, err := EmbedText("hello world!")
	if err != nil {
		t.Errorf("Expecting no error, got %s", err.Error())
	}
	if len(embd) != 768 {
		t.Errorf("Expecting the embedding to be of length 768, got %d", len(embd))
	}
}
