package goliblocalise

import (
	"context"
	"encoding/json"
	"io"
	"os"
)

type LocaliseService interface {
	ListLocaliseTranslation(ctx context.Context, req *ListLocaliseTranslationRequest) ([]LocaliseError, error)
}

type localiseService struct {
	ErrorFileName string
}

type ListLocaliseTranslationRequest struct {
	FileName string // FileName from which the data to be loaded for translation
}

// ListLocaliseTranslation implements LocaliseService.
func (l *localiseService) ListLocaliseTranslation(ctx context.Context, req *ListLocaliseTranslationRequest) ([]LocaliseError, error) {
	response := []LocaliseError{}

	resp, err := os.Open(l.ErrorFileName)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewLocaliseService(ctx context.Context, repoUrl string) LocaliseService {
	return &localiseService{
		ErrorFileName: "errors.json",
	}
}

type LocaliseError struct {
	Language       string `json:"language"`
	Key            string `json:"key"`
	Value          string `json:"value"`
	HTTPStatusCode int    `json:"httpStatusCode"`
}
