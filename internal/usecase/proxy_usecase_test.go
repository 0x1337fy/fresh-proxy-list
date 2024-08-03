package usecase

import (
	"errors"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/internal/service"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	testIP           = "13.37.0.1"
	testPort         = "8080"
	testProxy        = testIP + ":" + testPort
	testHTTPCategory = "HTTP"
	mockProxy        = entity.Proxy{
		Proxy:     testProxy,
		IP:        testIP,
		Port:      testPort,
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
)

func TestNewProxyUsecase(t *testing.T) {
	mockRepo := &mockProxyRepository{}
	mockProxyService := &mockProxyService{}
	usecase := NewProxyUsecase(mockRepo, mockProxyService)

	if _, ok := usecase.(*ProxyUsecase); !ok {
		t.Errorf("NewProxyUsecase() did not return a *ProxyUsecase")
	}

	pu, ok := usecase.(*ProxyUsecase)
	if !ok {
		t.Fatalf("Failed to cast to *ProxyUsecase")
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
			proxyService    service.ProxyServiceInterface
		}
		args struct {
			source entity.Source
			proxy  string
		}
		wantErr error
	}{
		{
			name: "Proxy Not Found",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
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
			name: "Proxy Format Incorrect",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "invalid-proxy",
			},
			wantErr: errors.New("proxy format incorrect"),
		},
		{
			name: "Proxy Format Not Match",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "invalid-proxy:1337",
			},
			wantErr: errors.New("proxy format not match"),
		},
		{
			name: "Proxy is Special IP With Exception IP",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "127.0.0.1:1337",
			},
			wantErr: errors.New("proxy belongs to special ip"),
		},
		{
			name: "Proxy is Special IP",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: "192.168.0.0:1337",
			},
			wantErr: errors.New("proxy belongs to special ip"),
		},
		{
			name: "Proxy port more than 65535",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  testHTTPCategory,
					IsChecked: false,
				},
				proxy: testIP + ":65540",
			},
			wantErr: errors.New("proxy port format incorrect"),
		},
		{
			name: "Proxy Has Been Processed",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
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
		{
			name: "Valid Proxy",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  mockProxy.Category,
					IsChecked: false,
				},
				proxy: mockProxy.Proxy,
			},
			wantErr: nil,
		},
		{
			name: "Not Valid Proxy",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService: &mockProxyService{
					CheckFunc: func(category string, ip string, port string) (*entity.Proxy, error) {
						return nil, errors.New("proxy not valid")
					},
				},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  mockProxy.Category,
					IsChecked: true,
				},
				proxy: mockProxy.Proxy,
			},
			wantErr: errors.New("proxy not valid"),
		},
		{
			name: "Valid Proxy with not checked",
			fields: struct {
				proxyRepository repository.ProxyRepositoryInterface
				proxyService    service.ProxyServiceInterface
			}{
				proxyRepository: &mockProxyRepository{},
				proxyService: &mockProxyService{
					CheckFunc: func(category string, ip string, port string) (*entity.Proxy, error) {
						return &mockProxy, nil
					},
				},
			},
			args: struct {
				source entity.Source
				proxy  string
			}{
				source: entity.Source{
					Category:  mockProxy.Category,
					IsChecked: false,
				},
				proxy: mockProxy.Proxy,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{
				proxyRepository: tt.fields.proxyRepository,
				proxyService:    tt.fields.proxyService,
				proxyMap:        sync.Map{},
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

func TestIsInvalidIP(t *testing.T) {
	type args struct {
		ip string
	}

	tests := []struct {
		name string
		args struct {
			ip string
		}
		want bool
	}{
		{
			name: "Test 255.255.255.255",
			args: args{
				ip: "255.255.255.255",
			},
			want: true,
		},
		{
			name: "Test ::1",
			args: args{
				ip: "::1",
			},
			want: false,
		},
		{
			name: "Test 300.300.300.300",
			args: args{
				ip: "300.300.300.300",
			},
			want: true,
		},
		{
			name: "Test 192.168.1",
			args: args{
				ip: "192.168.1",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{}
			got := uc.IsSpecialIP(tt.args.ip)
			if got != tt.want {
				t.Errorf("For IP '%s', want %v but got %v", tt.args.ip, tt.want, got)
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
