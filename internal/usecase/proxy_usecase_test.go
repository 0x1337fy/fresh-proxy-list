package usecase

import (
	"errors"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/internal/service"
	"net"
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
	specialIPs = []string{"1.1.1.1"}
	privateIPs = []net.IPNet{
		{IP: net.IP{2, 2, 2, 2}, Mask: net.CIDRMask(8, 32)},
		{IP: net.IP{3, 3, 3, 3}, Mask: net.CIDRMask(12, 32)},
		{IP: net.IP{4, 4, 4, 4}, Mask: net.CIDRMask(16, 32)},
		{IP: net.IP{5, 5, 5, 5}, Mask: net.CIDRMask(16, 32)},
	}
)

func TestNewProxyUsecase(t *testing.T) {
	mockRepo := &mockProxyRepository{}
	mockProxyService := &mockProxyService{}
	specialIPs := specialIPs
	privateIPs := privateIPs
	usecase := NewProxyUsecase(mockRepo, mockProxyService, specialIPs, privateIPs)

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

	if !reflect.DeepEqual(pu.specialIPs, specialIPs) {
		t.Errorf("Expected specialIPs %v, but got %v", specialIPs, pu.specialIPs)
	}

	if !reflect.DeepEqual(pu.privateIPs, privateIPs) {
		t.Errorf("Expected privateIPs %v, but got %v", privateIPs, pu.privateIPs)
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
			category  string
			proxy     string
			isChecked bool
		}
		want    *entity.Proxy
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     "   ",
				isChecked: false,
			},
			want:    nil,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     "invalid-proxy",
				isChecked: false,
			},
			want:    nil,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     "invalid-proxy:1337",
				isChecked: false,
			},
			want:    nil,
			wantErr: errors.New("proxy format not match"),
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     "1.1.1.1:1337",
				isChecked: false,
			},
			want:    nil,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     testIP + ":65540",
				isChecked: false,
			},
			want:    nil,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  testHTTPCategory,
				proxy:     testProxy,
				isChecked: false,
			},
			want:    nil,
			wantErr: errors.New("proxy has been processed"),
		},
		{
			name: "Valid Proxy",
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  mockProxy.Category,
				proxy:     mockProxy.Proxy,
				isChecked: true,
			},
			want: &entity.Proxy{
				Category:  mockProxy.Category,
				Proxy:     mockProxy.Proxy,
				IP:        mockProxy.IP,
				Port:      mockProxy.Port,
				TimeTaken: mockProxy.TimeTaken,
				CheckedAt: mockProxy.CheckedAt,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  mockProxy.Category,
				proxy:     mockProxy.Proxy,
				isChecked: true,
			},
			want:    nil,
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
				category  string
				proxy     string
				isChecked bool
			}{
				category:  mockProxy.Category,
				proxy:     mockProxy.Proxy,
				isChecked: false,
			},
			want: &entity.Proxy{
				Category:  mockProxy.Category,
				Proxy:     mockProxy.Proxy,
				IP:        mockProxy.IP,
				Port:      mockProxy.Port,
				TimeTaken: 0,
				CheckedAt: "",
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
				specialIPs:      specialIPs,
				privateIPs:      privateIPs,
			}

			if tt.name == "Proxy Has Been Processed" {
				uc.proxyMap.Store(tt.args.category+"_"+tt.args.proxy, true)
			}

			got, err := uc.ProcessProxy(tt.args.category, tt.args.proxy, tt.args.isChecked)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ProcessProxy() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSpecialIP(t *testing.T) {
	type args struct {
		ip string
	}

	tests := []struct {
		name   string
		fields struct {
			specialIPs []string
			privateIPs []net.IPNet
		}
		args struct {
			ip string
		}
		want bool
	}{
		{
			name: "Test 1.1.1.1",
			fields: struct {
				specialIPs []string
				privateIPs []net.IPNet
			}{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			},
			args: args{
				ip: "1.1.1.1",
			},
			want: true,
		},
		{
			name: "Test ::1",
			fields: struct {
				specialIPs []string
				privateIPs []net.IPNet
			}{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			},
			args: args{
				ip: "::1",
			},
			want: true,
		},
		{
			name: "Test 2.2.2.2",
			fields: struct {
				specialIPs []string
				privateIPs []net.IPNet
			}{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			},
			args: args{
				ip: "2.2.2.2",
			},
			want: true,
		},
		{
			name: "Test 192.168.1",
			fields: struct {
				specialIPs []string
				privateIPs []net.IPNet
			}{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			},
			args: args{
				ip: "192.168.1",
			},
			want: true,
		},
		{
			name: "Test 13.37.0.1",
			fields: struct {
				specialIPs []string
				privateIPs []net.IPNet
			}{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			},
			args: args{
				ip: "13.37.0.1",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProxyUsecase{
				specialIPs: specialIPs,
				privateIPs: privateIPs,
			}
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
