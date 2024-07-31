package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

const (
	methodGet       = "GET"
	errorStatusText = "failed to fetch data: %s"
	testSite        = "http://example.com"
	errNone         = "Expected error but got none"
	errNilData      = "Expected nil data but got %v"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type mockReadCloser struct {
	data []byte
	err  error
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	copy(p, m.data)
	return len(m.data), io.EOF
}

func (m *mockReadCloser) Close() error {
	return nil
}

type mockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestNewFetcher(t *testing.T) {
	client := http.DefaultTransport
	newFetcher := NewFetcher(client, func(method, url string, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, url, body)
	})
	if newFetcher == nil {
		t.Fatal("Expected non-nil FetcherInterface from NewFetcher")
	}

	fetcherUtilInstance, ok := newFetcher.(*FetcherUtil)
	if !ok {
		t.Fatal("Type assertion to *FetcherUtil failed")
	}

	req, err := fetcherUtilInstance.newRequest("GET", testSite, nil)
	if err != nil {
		t.Fatalf("Expected no error from newRequest but got: %v", err)
	}

	if req.Method != "GET" {
		t.Fatalf("Expected method to be GET, but got: %s", req.Method)
	}
	if req.URL.String() != testSite {
		t.Fatalf("Expected URL to be %s, but got: %s", testSite, req.URL.String())
	}

	if fetcherUtilInstance.client != http.DefaultTransport {
		t.Fatalf("Expected client to be http.DefaultTransport, but got: %v", fetcherUtilInstance.client)
	}
}

func TestFetchData(t *testing.T) {
	type fields struct {
		client     *mockTransport
		newRequest func(method string, url string, body io.Reader) (*http.Request, error)
	}

	type args struct {
		url string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				client: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(&mockReadCloser{data: []byte("response data")}),
					},
					err: nil,
				},
				newRequest: http.NewRequest,
			},
			args:    args{url: testSite},
			want:    []byte("response data"),
			wantErr: nil,
		},
		{
			name: "Request Error",
			fields: fields{
				client: &mockTransport{
					response: nil,
					err:      fmt.Errorf("request error"),
				},
				newRequest: http.NewRequest,
			},
			args:    args{url: testSite},
			want:    nil,
			wantErr: errors.New("request error"),
		},
		{
			name: "Response Error",
			fields: fields{
				client: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(&mockReadCloser{data: []byte("error response")}),
					},
					err: nil,
				},
				newRequest: http.NewRequest,
			},
			args:    args{url: testSite},
			want:    []byte("error response"),
			wantErr: errors.New("failed to fetch data: Internal Server Error"),
		},
		{
			name: "Read Body Error",
			fields: fields{
				client: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusOK,
						Body:       &mockReadCloser{err: fmt.Errorf("body read error")},
					},
					err: nil,
				},
				newRequest: http.NewRequest,
			},
			args:    args{url: testSite},
			want:    nil,
			wantErr: errors.New("body read error"),
		},
		{
			name: "New Request Error",
			fields: fields{
				client: &mockTransport{
					response: nil,
					err:      nil,
				},
				newRequest: func(method string, url string, body io.Reader) (*http.Request, error) {
					return nil, fmt.Errorf("new request error")
				},
			},
			args:    args{url: testSite},
			want:    nil,
			wantErr: errors.New("new request error"),
		},
		{
			name: "Do Error",
			fields: fields{
				client: &mockTransport{
					response: nil,
					err:      fmt.Errorf("do error"),
				},
				newRequest: func(method string, url string, body io.Reader) (*http.Request, error) {
					return &http.Request{}, nil
				},
			},
			args:    args{url: testSite},
			want:    nil,
			wantErr: errors.New("do error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcherUtil := &FetcherUtil{
				client:     tt.fields.client,
				newRequest: tt.fields.newRequest,
			}

			got, err := fetcherUtil.FetchData(tt.args.url)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("FetchData() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("FetchData() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		body   io.Reader
	}
	type wantData struct {
		url    string
		method string
	}

	mockClient := &mockRoundTripper{}
	f := NewFetcher(mockClient, http.NewRequest)

	tests := []struct {
		name     string
		args     args
		wantData wantData
		wantErr  error
	}{
		{
			name: "GET Request",
			args: args{
				method: http.MethodGet,
				url:    testSite,
				body:   nil,
			},
			wantData: wantData{
				url:    testSite,
				method: http.MethodGet,
			},
			wantErr: nil,
		},
		{
			name: "POST Request",
			args: args{
				method: http.MethodPost,
				url:    testSite,
				body:   bytes.NewReader([]byte("body")),
			},
			wantData: wantData{
				url:    testSite,
				method: http.MethodPost,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := f.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if req.Method != tt.wantData.method {
				t.Fatalf("expected method: %s, got: %s", tt.wantData.method, req.Method)
			}

			if req.URL.String() != tt.wantData.url {
				t.Fatalf("expected URL: %s, got: %s", tt.wantData.url, req.URL.String())
			}
		})
	}
}
