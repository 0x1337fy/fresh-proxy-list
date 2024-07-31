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
	testIP           = "192.168.1.1"
	testPort         = "8080"
	testProxy        = "192.168.1.1:8080"
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
				proxy: "invalidProxy",
			},
			wantErr: errors.New("proxy format incorrect"),
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
