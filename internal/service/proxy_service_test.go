package service

import (
	"errors"
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/pkg/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	testIP             = "192.168.1.1"
	testPort           = "8080"
	testProxy          = "192.168.1.1:8080"
	testHTTPCategory   = "HTTP"
	testHTTPSCategory  = "HTTPS"
	testSOCKS5Category = "SOCKS5"
	proxyService       = &ProxyService{
		httpTestingSites:  []string{"http://test1.com", "http://test2.com"},
		httpsTestingSites: []string{"http://secure1.com", "http://secure2.com"},
		userAgents:        []string{"Mozilla", "Chrome", "Safari"},
	}
)

type mockURLParserUtil struct {
	ParseFunc func(urlStr string) (*url.URL, error)
}

func (m *mockURLParserUtil) Parse(urlStr string) (*url.URL, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(urlStr)
	}
	return url.Parse(urlStr)
}

type mockFetcherUtil struct {
	fetchDataByte  []byte
	fetcherError   error
	NewRequestFunc func(method, url string, body io.Reader) (*http.Request, error)
	DoFunc         func(client *http.Client, req *http.Request) (*http.Response, error)
}

func (m *mockFetcherUtil) FetchData(url string) ([]byte, error) {
	if m.fetcherError != nil {
		return nil, m.fetcherError
	}
	return m.fetchDataByte, nil
}

func (m *mockFetcherUtil) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(client, req)
	}
	return httptest.NewRecorder().Result(), nil
}

func (m *mockFetcherUtil) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if m.NewRequestFunc != nil {
		return m.NewRequestFunc(method, url, body)
	}
	return http.NewRequest(method, url, body)
}

func TestNewProxyService(t *testing.T) {
	mockFetcher := &mockFetcherUtil{}
	mockURLParser := &mockURLParserUtil{}
	httpTestingSites := proxyService.httpTestingSites
	httpsTestingSites := proxyService.httpsTestingSites
	userAgents := proxyService.userAgents
	service := NewProxyService(mockFetcher, mockURLParser, httpTestingSites, httpsTestingSites, userAgents)

	if _, ok := service.(*ProxyService); !ok {
		t.Errorf("NewProxyService() did not return a *ProxyService")
	}

	pu, ok := service.(*ProxyService)
	if !ok {
		t.Fatalf("Failed to cast to *ProxyService")
	}

	if !reflect.DeepEqual(pu.httpTestingSites, proxyService.httpTestingSites) {
		t.Errorf("Expected httpTestingSites %v, but got %v", proxyService.httpTestingSites, pu.httpTestingSites)
	}

	if !reflect.DeepEqual(pu.httpsTestingSites, proxyService.httpsTestingSites) {
		t.Errorf("Expected httpsTestingSites %v, but got %v", proxyService.httpsTestingSites, pu.httpsTestingSites)
	}

	if !reflect.DeepEqual(pu.userAgents, proxyService.userAgents) {
		t.Errorf("Expected userAgents %v, but got %v", proxyService.userAgents, pu.userAgents)
	}
}

func TestCheck(t *testing.T) {
	type fields struct {
		fetcherUtil   utils.FetcherUtilInterface
		urlParserUtil utils.URLParserUtilInterface
	}

	type args struct {
		category string
		ip       string
		port     string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Proxy
		wantErr error
	}{
		{
			name: "Test valid",
			fields: fields{
				fetcherUtil:   &mockFetcherUtil{},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPSCategory,
				ip:       testIP,
				port:     testPort,
			},
			want: &entity.Proxy{
				Category:  testHTTPSCategory,
				Proxy:     testProxy,
				IP:        testIP,
				Port:      testPort,
				TimeTaken: 123.45,
				CheckedAt: time.Now().Format(time.RFC3339),
			},
			wantErr: nil,
		},
		{
			name: "Test error parse url",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{},
				urlParserUtil: &mockURLParserUtil{
					ParseFunc: func(urlStr string) (*url.URL, error) {
						return nil, errors.New("parse error")
					},
				},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:    nil,
			wantErr: errors.New("error parsing proxy URL: parse error"),
		},
		{
			name: "Test creating request",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					NewRequestFunc: func(method, url string, body io.Reader) (*http.Request, error) {
						return nil, errors.New("error creating request")
					},
				},
			},
			args: args{
				category: testSOCKS5Category,
				ip:       testIP,
				port:     testPort,
			},
			want:    nil,
			wantErr: errors.New("error creating request: error creating request"),
		},
		{
			name:   "Test unsupported proxy category",
			fields: fields{},
			args: args{
				category: "FTP",
				ip:       testIP,
				port:     testPort,
			},
			want:    nil,
			wantErr: errors.New("proxy category FTP not supported"),
		},
		{
			name: "Test request error",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client *http.Client, req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("network error")
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:    nil,
			wantErr: errors.New("request error: network error"),
		},
		{
			name: "Test unexpected status code",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client *http.Client, req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       http.NoBody,
						}, nil
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:    nil,
			wantErr: errors.New("unexpected status code 500: Internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				fetcherUtil:       tt.fields.fetcherUtil,
				urlParserUtil:     tt.fields.urlParserUtil,
				httpTestingSites:  proxyService.httpTestingSites,
				httpsTestingSites: proxyService.httpsTestingSites,
				userAgents:        proxyService.userAgents,
				semaphore:         make(chan struct{}, 10),
			}

			got, err := s.Check(tt.args.category, tt.args.ip, tt.args.port)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("ProxyService.Check() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("ProxyService.Check() error = %v, want %v", err, tt.wantErr)
			}

			if tt.want != nil && (!reflect.DeepEqual(got.Category, tt.want.Category) ||
				!reflect.DeepEqual(got.Proxy, tt.want.Proxy) ||
				!reflect.DeepEqual(got.IP, tt.want.IP) ||
				!reflect.DeepEqual(got.Port, tt.want.Port)) {
				t.Errorf("ProxyService.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestingSite(t *testing.T) {
	type fields struct {
		httpTestingSites  []string
		httpsTestingSites []string
	}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "HTTP",
			fields: fields{
				httpTestingSites: proxyService.httpTestingSites,
			},
			want: proxyService.httpTestingSites,
		},
		{
			name: "HTTPS",
			fields: fields{
				httpsTestingSites: proxyService.httpsTestingSites,
			},
			want: proxyService.httpsTestingSites,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				httpTestingSites:  tt.fields.httpTestingSites,
				httpsTestingSites: tt.fields.httpsTestingSites,
			}

			site := s.GetTestingSite(tt.name)
			if len(site) == 0 {
				t.Errorf("expected a non-empty site for name %s", tt.name)
			}

			found := false
			for _, expectedSite := range tt.want {
				if expectedSite == site {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected site to be from %s sites, got %s, wants: %v", tt.name, site, tt.want)
			}
		})
	}
}

func TestGetRandomUserAgent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Random user agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				userAgents: proxyService.userAgents,
			}

			site := s.GetRandomUserAgent()
			found := false
			for _, ua := range s.userAgents {
				if ua == site {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected user agent not found, got %s, want: %v", site, s.userAgents)
			}
		})
	}
}
