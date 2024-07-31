package usecase

import (
	"fresh-proxy-list/internal/entity"
	"sync"
)

type mockProxyService struct {
	CheckFunc              func(category string, ip string, port string) (entity.Proxy, error)
	GetTestingSiteFunc     func(category string) string
	GetRandomUserAgentFunc func() string
}

func (m *mockProxyService) Check(category string, ip string, port string) (entity.Proxy, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(category, ip, port)
	}
	return entity.Proxy{}, nil
}

func (m *mockProxyService) GetTestingSite(category string) string {
	if m.GetTestingSiteFunc != nil {
		return m.GetTestingSiteFunc(category)
	}
	return ""
}

func (m *mockProxyService) GetRandomUserAgent() string {
	if m.GetRandomUserAgentFunc != nil {
		return m.GetRandomUserAgentFunc()
	}
	return ""
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
