package rag

import (
	"context"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
)

type EmbeddingModel struct {
	model textencoding.Interface
}

func NewEmbeddingModel() *EmbeddingModel {
	model, err := tasks.Load[textencoding.Interface](&tasks.Config{
		ModelsDir: "./embedding_cache",
		ModelName: "sentence-transformers/all-MiniLM-L6-v2",
	})
	if err != nil {
		panic(err)
	}
	return &EmbeddingModel{model: model}
}

func (e *EmbeddingModel) GetEmbedding(text string) ([]float32, error) {
	result, err := e.model.Encode(context.Background(), text, int(bert.MeanPooling))
	if err != nil {
		return nil, err
	}
	return result.Vector.Data().F32(), nil
}
