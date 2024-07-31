package repository

import (
	"errors"
	"fresh-proxy-list/internal/entity"
	"reflect"
	"testing"
)

func createMockRepository(sourcesJSON string) SourceRepositoryInterface {
	return NewSourceRepository(sourcesJSON)
}

func TestLoadSources(t *testing.T) {
	type args struct {
		proxy_resources string
	}

	tests := []struct {
		name    string
		args    args
		want    []entity.Source
		wantErr error
	}{
		{
			name: "Empty resources",
			args: args{
				proxy_resources: "",
			},
			want:    nil,
			wantErr: errors.New("PROXY_RESOURCES not found on environment"),
		},
		{
			name: "Invalid JSON",
			args: args{
				proxy_resources: `{"invalid": "json"`,
			},
			want:    nil,
			wantErr: errors.New("error parsing JSON: unexpected end of JSON input"),
		},
		{
			name: "Valid JSON",
			args: args{
				proxy_resources: `[{"method": "GET", "category": "general", "url": "http://example.com", "is_checked": true}]`,
			},
			want: []entity.Source{
				{
					Method:    "GET",
					Category:  "general",
					URL:       "http://example.com",
					IsChecked: true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createMockRepository(tt.args.proxy_resources)
			sources, err := repo.LoadSources()
			if !reflect.DeepEqual(sources, tt.want) {
				t.Errorf("LoadSources() = %v, want %v", sources, tt.want)
			}
			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("LoadSources() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
