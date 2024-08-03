package usecase

import (
	"fresh-proxy-list/internal/entity"
	"path/filepath"
	"strings"
	"testing"
)

const (
	storageDir    = "storage"
	classicDir    = "classic"
	advancedDir   = "advanced"
	txtExtension  = "txt"
	csvExtension  = "csv"
	jsonExtension = "json"
	xmlExtension  = "xml"
	yamlExtension = "yaml"
)

func TestSaveFiles(t *testing.T) {
	mockFileRepo := &mockFileRepository{}
	mockProxyRepo := &mockProxyRepository{}
	fileOutputExtensions := []string{"txt", "csv"}
	usecase := NewFileUsecase(mockFileRepo, mockProxyRepo, fileOutputExtensions)

	mockProxyRepo.GetAllClassicViewFunc = func() []string {
		return []string{"proxy1", "proxy2"}
	}
	mockProxyRepo.GetAllAdvancedViewFunc = func() []entity.AdvancedProxy {
		return []entity.AdvancedProxy{
			{Proxy: "proxy1", IP: "1.1.1.1", Port: "80", Categories: []string{"category1"}},
			{Proxy: "proxy2", IP: "2.2.2.2", Port: "8080", Categories: []string{"category2"}},
		}
	}
	mockProxyRepo.GetHTTPClassicViewFunc = func() []string {
		return []string{"httpProxy1", "httpProxy2"}
	}
	mockProxyRepo.GetHTTPAdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			{Category: "httpCategory1", Proxy: "httpProxy1", IP: "3.3.3.3", Port: "8081"},
			{Category: "httpCategory2", Proxy: "httpProxy2", IP: "4.4.4.4", Port: "8082"},
		}
	}
	mockProxyRepo.GetHTTPSClassicViewFunc = func() []string {
		return []string{"httpsProxy1", "httpsProxy2"}
	}
	mockProxyRepo.GetHTTPSAdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			{Category: "httpsCategory1", Proxy: "httpsProxy1", IP: "5.5.5.5", Port: "8083"},
			{Category: "httpsCategory2", Proxy: "httpsProxy2", IP: "6.6.6.6", Port: "8084"},
		}
	}
	mockProxyRepo.GetSOCKS4ClassicViewFunc = func() []string {
		return []string{"socks4Proxy1", "socks4Proxy2"}
	}
	mockProxyRepo.GetSOCKS4AdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			{Category: "socks4Category1", Proxy: "socks4Proxy1", IP: "7.7.7.7", Port: "8085"},
			{Category: "socks4Category2", Proxy: "socks4Proxy2", IP: "8.8.8.8", Port: "8086"},
		}
	}
	mockProxyRepo.GetSOCKS5ClassicViewFunc = func() []string {
		return []string{"socks5Proxy1", "socks5Proxy2"}
	}
	mockProxyRepo.GetSOCKS5AdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			{Category: "socks5Category1", Proxy: "socks5Proxy1", IP: "9.9.9.9", Port: "8087"},
			{Category: "socks5Category2", Proxy: "socks5Proxy2", IP: "10.10.10.10", Port: "8088"},
		}
	}

	calls := 0
	mockFileRepo.SaveFileFunc = func(filename string, data interface{}, ext string) error {
		calls++
		t.Logf("SaveFile called with filename: %s, extension: %s", filename, ext)
		if !strings.HasPrefix(filename, filepath.Join(storageDir, classicDir)) &&
			!strings.HasPrefix(filename, filepath.Join(storageDir, advancedDir)) {
			t.Errorf("Unexpected filename: %s", filename)
		}
		if ext != csvExtension && ext != jsonExtension && ext != xmlExtension && ext != yamlExtension && ext != txtExtension {
			t.Errorf("Unexpected extension: %s", ext)
		}
		return nil
	}

	usecase.SaveFiles()

	// (5 categories * 2 extensions * 2 file types (classic, advanced)) + (5 all * 1 extension txt * 1 file type classic)
	expectedCalls := (5 * 2 * 2) + (5 * 1 * 1)
	if calls != expectedCalls {
		t.Errorf("Expected %d calls, but got %d", expectedCalls, calls)
	}
}
