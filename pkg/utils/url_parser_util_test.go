package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestUtilParse(t *testing.T) {
	type args struct {
		rawURL string
	}

	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr error
	}{
		{
			name: "Valid URL",
			args: args{
				rawURL: "http://example.com",
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
			wantErr: nil,
		},
		{
			name: "URL with path",
			args: args{
				rawURL: "https://example.com/path?query=1",
			},
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/path",
				RawQuery: "query=1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLParser()

			got, err := u.Parse(tt.args.rawURL)
			if tt.wantErr != nil {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
