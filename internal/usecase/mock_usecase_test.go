package usecase

import (
	"fresh-proxy-list/internal/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
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

func (m *mockFetcherUtil) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	if m.NewRequestFunc != nil {
		return m.NewRequestFunc(method, url, body)
	}
	return http.NewRequest(method, url, body)
}

type mockSourceRepository struct {
	LoadSourcesFunc func() ([]entity.Source, error)
}

func (m *mockSourceRepository) LoadSources() ([]entity.Source, error) {
	return m.LoadSourcesFunc()
}

type mockFileRepository struct {
	SaveFileFunc func(filename string, data interface{}, ext string) error
}

func (m *mockFileRepository) SaveFile(filename string, data interface{}, ext string) error {
	if m.SaveFileFunc != nil {
		return m.SaveFileFunc(filename, data, ext)
	}
	return nil
}

type mockProxyRepository struct {
	StoreFunc                 func(proxy entity.Proxy)
	GetAllClassicViewFunc     func() []string
	GetHTTPClassicViewFunc    func() []string
	GetHTTPSClassicViewFunc   func() []string
	GetSOCKS4ClassicViewFunc  func() []string
	GetSOCKS5ClassicViewFunc  func() []string
	GetAllAdvancedViewFunc    func() []entity.AdvancedProxy
	GetHTTPAdvancedViewFunc   func() []entity.Proxy
	GetHTTPSAdvancedViewFunc  func() []entity.Proxy
	GetSOCKS4AdvancedViewFunc func() []entity.Proxy
	GetSOCKS5AdvancedViewFunc func() []entity.Proxy

	storedProxies      []entity.Proxy
	mu                 sync.Mutex
	IsProxyWorkingFunc func(entity.Source, string, string) (entity.Proxy, error)
}

func (m *mockProxyRepository) Store(proxy entity.Proxy) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storedProxies = append(m.storedProxies, proxy)
}

func (m *mockProxyRepository) GetStoredProxies() []entity.Proxy {
	return m.storedProxies
}

func (m *mockProxyRepository) GetAllClassicView() []string {
	if m.GetAllClassicViewFunc != nil {
		return m.GetAllClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPClassicView() []string {
	if m.GetHTTPClassicViewFunc != nil {
		return m.GetHTTPClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPSClassicView() []string {
	if m.GetHTTPSClassicViewFunc != nil {
		return m.GetHTTPSClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS4ClassicView() []string {
	if m.GetSOCKS4ClassicViewFunc != nil {
		return m.GetSOCKS4ClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS5ClassicView() []string {
	if m.GetSOCKS5ClassicViewFunc != nil {
		return m.GetSOCKS5ClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	if m.GetAllAdvancedViewFunc != nil {
		return m.GetAllAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	if m.GetHTTPAdvancedViewFunc != nil {
		return m.GetHTTPAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	if m.GetHTTPSAdvancedViewFunc != nil {
		return m.GetHTTPSAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	if m.GetSOCKS4AdvancedViewFunc != nil {
		return m.GetSOCKS4AdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	if m.GetSOCKS5AdvancedViewFunc != nil {
		return m.GetSOCKS5AdvancedViewFunc()
	}
	return nil
}
