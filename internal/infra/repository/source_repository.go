package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"fresh-proxy-list/internal/entity"
)

type sourceRepository struct {
	proxyResources string
}

type SourceRepositoryInterface interface {
	LoadSources() ([]entity.Source, error)
}

func NewSourceRepository(proxyResources string) SourceRepositoryInterface {
	return &sourceRepository{
		proxyResources: proxyResources,
	}
}

func (r *sourceRepository) LoadSources() ([]entity.Source, error) {
	sourcesJSON := r.proxyResources
	if sourcesJSON == "" {
		return nil, errors.New("PROXY_RESOURCES not found on environment")
	}

	var sources []entity.Source
	err := json.Unmarshal([]byte(sourcesJSON), &sources)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return sources, nil
}
