package usecase

import (
	"errors"
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/config"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/pkg/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
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
	mockProxy          = entity.Proxy{
		Proxy:     testProxy,
		IP:        "192.168.1.1",
		Port:      "8080",
		Category:  testHTTPCategory,
		TimeTaken: 0,
		CheckedAt: time.Now().Format(time.RFC3339),
	}
	mockAdvancedProxy = entity.AdvancedProxy{
		Proxy:      mockProxy.Proxy,
		IP:         mockProxy.IP,
		Port:       mockProxy.Port,
		TimeTaken:  mockProxy.TimeTaken,
		CheckedAt:  mockProxy.CheckedAt,
		Categories: []string{testHTTPCategory},
	}
	proxyUsecase = &ProxyUsecase{
		httpTestingSites:  []string{"http://test1.com", "http://test2.com"},
		httpsTestingSites: []string{"http://secure1.com", "http://secure2.com"},
		userAgents:        []string{"Mozilla", "Chrome", "Safari"},
	}
)

func TestNewProxyUsecase(t *testing.T) {
	mockRepo := &mockProxyRepository{}
	mockFetcher := &mockFetcherUtil{}
	mockURLParser := &mockURLParserUtil{}

	usecase := NewProxyUsecase(mockRepo, mockFetcher, mockURLParser)

	if _, ok := usecase.(*ProxyUsecase); !ok {
		t.Errorf("NewProxyUsecase() did not return a *ProxyUsecase")
	}

	pu, ok := usecase.(*ProxyUsecase)
	if !ok {
		t.Fatalf("Failed to cast to *ProxyUsecase")
	}

	if !reflect.DeepEqual(pu.httpTestingSites, config.HTTPTestingSites) {
		t.Errorf("Expected httpTestingSites %v, but got %v", config.HTTPTestingSites, pu.httpTestingSites)
	}

	if !reflect.DeepEqual(pu.httpsTestingSites, config.HTTPSTestingSites) {
		t.Errorf("Expected httpsTestingSites %v, but got %v", config.HTTPSTestingSites, pu.httpsTestingSites)
	}

	if !reflect.DeepEqual(pu.userAgents, config.UserAgents) {
		t.Errorf("Expected userAgents %v, but got %v", config.UserAgents, pu.userAgents)
	}

	testKey := testHTTPCategory + "_" + testProxy
	testValue := true
	pu.proxyMap.Store(testKey, testValue)

	value, ok := pu.proxyMap.Load(testKey)
	if !ok || value != testValue {
		t.Errorf("Expected value %v, but got %v", testValue, value)
	}

	_, loaded := pu.proxyMap.LoadOrStore(testKey, false)
	if !loaded {
		t.Errorf("Expected LoadOrStore to return true indicating the key was loaded")
	}

	value, _ = pu.proxyMap.Load(testKey)
	if value != testValue {
		t.Errorf("Expected value %v after LoadOrStore, but got %v", testValue, value)
	}
}

func TestProcessProxy(t *testing.T) {
	tests := []struct {
		name   string
		fields struct {
			proxyRepository repository.ProxyRepositoryInterface
			fetcherUtil     utils.FetcherUtilInterface
			urlParserUtil   utils.URLParserUtilInterface
		}
		args struct {
			source entity.Source
			proxy  string
		}
		wantErr error
	}{
		// {
		// 	name: "Valid Proxy",
		// 	fields: struct {
		// 		proxyRepository repository.ProxyRepositoryInterface
		// 		fetcherUtil     utils.FetcherUtilInterface
		// 		urlParserUtil   utils.URLParserUtilInterface
		// 	}{
		// 		proxyRepository: &mockProxyRepository{},
		// 		fetcherUtil: &mockFetcherUtil{
		// 			NewRequestFunc: func(method, url string, body io.Reader) (*http.Request, error) {
		// 				return httptest.NewRequest(method, url, body), nil
		// 			},
		// 		},
		// 		urlParserUtil: &mockURLParserUtil{},
		// 	},
		// 	args: struct {
		// 		source entity.Source
		// 		proxy  string
		// 	}{
		// 		source: entity.Source{
		// 			Category:  testHTTPCategory,
		// 			IsChecked: true,
		// 		},
		// 		proxy: testProxy,
		// 	},
		// 	wantErr: nil,
		// },
		{
			name: "Valid Proxy with not checked",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				fetcherUtil     utils.FetcherUtilInterface
				urlParserUtil   utils.URLParserUtilInterface
			}{
				proxyRepository: &mockProxyRepository{},
				fetcherUtil: &mockFetcherUtil{
					NewRequestFunc: func(method, url string, body io.Reader) (*http.Request, error) {
						return httptest.NewRequest(method, url, body), nil
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: testProxy,
			},
			wantErr: nil,
		},
		{
			name: "Proxy Format Incorrect",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				fetcherUtil     utils.FetcherUtilInterface
				urlParserUtil   utils.URLParserUtilInterface
			}{
				proxyRepository: &mockProxyRepository{},
				fetcherUtil:     &mockFetcherUtil{},
				urlParserUtil:   &mockURLParserUtil{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "invalidProxy",
			},
			wantErr: errors.New("proxy format incorrect"),
		},
		{
			name: "Proxy Not Found",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				fetcherUtil     utils.FetcherUtilInterface
				urlParserUtil   utils.URLParserUtilInterface
			}{
				proxyRepository: &mockProxyRepository{},
				fetcherUtil:     &mockFetcherUtil{},
				urlParserUtil:   &mockURLParserUtil{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "   ",
			},
			wantErr: errors.New("proxy not found"),
		},
		{
			name: "Proxy Has Been Processed",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				fetcherUtil     utils.FetcherUtilInterface
				urlParserUtil   utils.URLParserUtilInterface
			}{
				proxyRepository: &mockProxyRepository{},
				fetcherUtil: &mockFetcherUtil{
					NewRequestFunc: func(method, url string, body io.Reader) (*http.Request, error) {
						return httptest.NewRequest(method, url, body), nil
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: testProxy,
			},
			wantErr: errors.New("proxy has been processed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{
				proxyRepository:   tt.fields.proxyRepository,
				fetcherUtil:       tt.fields.fetcherUtil,
				urlParserUtil:     tt.fields.urlParserUtil,
				httpTestingSites:  config.HTTPTestingSites,
				httpsTestingSites: config.HTTPSTestingSites,
				userAgents:        config.UserAgents,
				proxyMap:          sync.Map{},
				semaphore:         make(chan struct{}, 100),
			}

			if tt.name == "Proxy Has Been Processed" {
				uc.proxyMap.Store(tt.args.source.Category+"_"+tt.args.proxy, true)
			}

			err := uc.ProcessProxy(tt.args.source, tt.args.proxy)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ProcessProxy() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsProxyWorking(t *testing.T) {
	type fields struct {
		fetcherUtil   utils.FetcherUtilInterface
		urlParserUtil utils.URLParserUtilInterface
	}

	type args struct {
		source entity.Source
		ip     string
		port   string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Proxy
		wantErr error
	}{
		{
			name: "Test valid",
			fields: fields{
				fetcherUtil:   &mockFetcherUtil{},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				source: entity.Source{
					Category: testHTTPSCategory,
				},
				ip:   testIP,
				port: testPort,
			},
			want: entity.Proxy{
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
				source: entity.Source{
					Category: testHTTPCategory,
				},
				ip:   testIP,
				port: testPort,
			},
			want:    entity.Proxy{},
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
				source: entity.Source{
					Category: testSOCKS5Category,
				},
				ip:   testIP,
				port: testPort,
			},
			want:    entity.Proxy{},
			wantErr: errors.New("error creating request: error creating request"),
		},
		{
			name:   "Test unsupported proxy category",
			fields: fields{},
			args: args{
				source: entity.Source{
					Category: "FTP",
				},
				ip:   testIP,
				port: testPort,
			},
			want:    entity.Proxy{},
			wantErr: errors.New("proxy category FTP not supported"),
		},
		{
			name: "Test request error",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client http.RoundTripper, req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("network error")
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				source: entity.Source{
					Category: testHTTPCategory,
				},
				ip:   testIP,
				port: testPort,
			},
			want:    entity.Proxy{},
			wantErr: errors.New("request error: network error"),
		},
		{
			name: "Test unexpected status code",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client http.RoundTripper, req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       http.NoBody,
						}, nil
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				source: entity.Source{
					Category: testHTTPCategory,
				},
				ip:   testIP,
				port: testPort,
			},
			want:    entity.Proxy{},
			wantErr: errors.New("unexpected status code 500: Internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{
				fetcherUtil:       tt.fields.fetcherUtil,
				urlParserUtil:     tt.fields.urlParserUtil,
				httpTestingSites:  proxyUsecase.httpTestingSites,
				httpsTestingSites: proxyUsecase.httpsTestingSites,
				userAgents:        proxyUsecase.userAgents,
				semaphore:         make(chan struct{}, 10),
			}

			got, err := uc.IsProxyWorking(tt.args.source, tt.args.ip, tt.args.port)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("ProxyUsecase.IsProxyWorking() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("ProxyUsecase.IsProxyWorking() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got.Category, tt.want.Category) ||
				!reflect.DeepEqual(got.Proxy, tt.want.Proxy) ||
				!reflect.DeepEqual(got.IP, tt.want.IP) ||
				!reflect.DeepEqual(got.Port, tt.want.Port) {
				t.Errorf("ProxyUsecase.IsProxyWorking() = %v, want %v", got, tt.want)
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
				httpTestingSites: proxyUsecase.httpTestingSites,
			},
			want: proxyUsecase.httpTestingSites,
		},
		{
			name: "HTTPS",
			fields: fields{
				httpsTestingSites: proxyUsecase.httpsTestingSites,
			},
			want: proxyUsecase.httpsTestingSites,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{
				httpTestingSites:  tt.fields.httpTestingSites,
				httpsTestingSites: tt.fields.httpsTestingSites,
			}

			site := uc.GetTestingSite(tt.name)
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
			uc := &ProxyUsecase{
				userAgents: proxyUsecase.userAgents,
			}

			site := uc.GetRandomUserAgent()
			found := false
			for _, ua := range uc.userAgents {
				if ua == site {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected user agent not found, got %s, want: %v", site, uc.userAgents)
			}
		})
	}
}

func TestGetAllAdvancedView(t *testing.T) {
	uc := &ProxyUsecase{
		proxyRepository: &mockProxyRepository{
			GetAllAdvancedViewFunc: func() []entity.AdvancedProxy {
				return []entity.AdvancedProxy{mockAdvancedProxy}
			},
		},
	}

	tests := []struct {
		name string
		want []entity.AdvancedProxy
	}{
		{
			name: "Should return all advanced view proxies",
			want: []entity.AdvancedProxy{mockAdvancedProxy},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uc.GetAllAdvancedView()
			if len(got) != len(tt.want) {
				t.Errorf("GetAllAdvancedView() = %v, want %v", got, tt.want)
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetAllAdvancedView() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
