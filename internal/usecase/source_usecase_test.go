package usecase

import (
	"errors"
	"fresh-proxy-list/internal/entity"
	"testing"
)

var (
	errMock = errors.New("mock error")
)

func TestLoadSourcesSuccess(t *testing.T) {
	mockRepo := &mockSourceRepository{
		LoadSourcesFunc: func() ([]entity.Source, error) {
			return []entity.Source{
				{Method: "GET", Category: "HTTP", URL: "http://example.com", IsChecked: true},
			}, nil
		},
	}

	uc := NewSourceUsecase(mockRepo)
	sources, err := uc.LoadSources()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expectedSources := []entity.Source{
		{Method: "GET", Category: "HTTP", URL: "http://example.com", IsChecked: true},
	}

	if len(sources) != len(expectedSources) {
		t.Errorf("expected %v sources, got %v", len(expectedSources), len(sources))
	}

	for i, source := range sources {
		if source != expectedSources[i] {
			t.Errorf("expected source %v, got %v", expectedSources[i], source)
		}
	}
}

func TestLoadSourcesError(t *testing.T) {
	mockRepo := &mockSourceRepository{
		LoadSourcesFunc: func() ([]entity.Source, error) {
			return nil, errMock
		},
	}

	uc := NewSourceUsecase(mockRepo)
	sources, err := uc.LoadSources()

	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if sources != nil {
		t.Errorf("expected nil sources, got %v", sources)
	}
}
