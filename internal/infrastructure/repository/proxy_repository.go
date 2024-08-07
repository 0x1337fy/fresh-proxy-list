package repository

import (
	"cmp"
	"fresh-proxy-list/internal/entity"
	"slices"
	"sync"
)

type ProxyRepository struct {
	Mutex              sync.RWMutex
	AllClassicView     []string
	HTTPClassicView    []string
	HTTPSClassicView   []string
	SOCKS4ClassicView  []string
	SOCKS5ClassicView  []string
	AllAdvancedView    []entity.AdvancedProxy
	HTTPAdvancedView   []entity.Proxy
	HTTPSAdvancedView  []entity.Proxy
	SOCKS4AdvancedView []entity.Proxy
	SOCKS5AdvancedView []entity.Proxy
}

type ProxyRepositoryInterface interface {
	Store(proxy *entity.Proxy)
	GetAllClassicView() []string
	GetHTTPClassicView() []string
	GetHTTPSClassicView() []string
	GetSOCKS4ClassicView() []string
	GetSOCKS5ClassicView() []string
	GetAllAdvancedView() []entity.AdvancedProxy
	GetHTTPAdvancedView() []entity.Proxy
	GetHTTPSAdvancedView() []entity.Proxy
	GetSOCKS4AdvancedView() []entity.Proxy
	GetSOCKS5AdvancedView() []entity.Proxy
}

func NewProxyRepository() ProxyRepositoryInterface {
	return &ProxyRepository{
		Mutex: sync.RWMutex{},
	}
}

func (r *ProxyRepository) Store(proxy *entity.Proxy) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	updateProxyAll := func(proxy *entity.Proxy, classicList *[]string, advancedList *[]entity.AdvancedProxy) {
		n, found := slices.BinarySearchFunc(*advancedList, entity.AdvancedProxy{Proxy: proxy.Proxy}, func(a, b entity.AdvancedProxy) int {
			return cmp.Compare(a.Proxy, b.Proxy)
		})
		if found {
			if proxy.Category == "HTTP" && proxy.TimeTaken > 0 {
				(*advancedList)[n].TimeTaken = proxy.TimeTaken
			}

			if m, found := slices.BinarySearch((*advancedList)[n].Categories, proxy.Category); !found {
				(*advancedList)[n].Categories = slices.Insert((*advancedList)[n].Categories, m, proxy.Category)
			}
		} else {
			*classicList = append(*classicList, proxy.Proxy)
			*advancedList = slices.Insert(*advancedList, n, entity.AdvancedProxy{
				Proxy:     proxy.Proxy,
				IP:        proxy.IP,
				Port:      proxy.Port,
				TimeTaken: proxy.TimeTaken,
				CheckedAt: proxy.CheckedAt,
				Categories: []string{
					proxy.Category,
				},
			})
		}
	}

	switch proxy.Category {
	case "HTTP":
		var (
			HTTPClassicView  = &r.HTTPClassicView
			HTTPAdvancedView = &r.HTTPAdvancedView
		)
		*HTTPClassicView = append(*HTTPClassicView, proxy.Proxy)
		*HTTPAdvancedView = append(*HTTPAdvancedView, *proxy)
	case "HTTPS":
		var (
			HTTPSClassicView  = &r.HTTPSClassicView
			HTTPSAdvancedView = &r.HTTPSAdvancedView
		)
		*HTTPSClassicView = append(*HTTPSClassicView, proxy.Proxy)
		*HTTPSAdvancedView = append(*HTTPSAdvancedView, *proxy)
	case "SOCKS4":
		var (
			SOCKS4ClassicView  = &r.SOCKS4ClassicView
			SOCKS4AdvancedView = &r.SOCKS4AdvancedView
		)
		*SOCKS4ClassicView = append(*SOCKS4ClassicView, proxy.Proxy)
		*SOCKS4AdvancedView = append(*SOCKS4AdvancedView, *proxy)
	case "SOCKS5":
		var (
			SOCKS5ClassicView  = &r.SOCKS5ClassicView
			SOCKS5AdvancedView = &r.SOCKS5AdvancedView
		)
		*SOCKS5ClassicView = append(*SOCKS5ClassicView, proxy.Proxy)
		*SOCKS5AdvancedView = append(*SOCKS5AdvancedView, *proxy)
	}

	updateProxyAll(proxy, &r.AllClassicView, &r.AllAdvancedView)
}

func (r *ProxyRepository) GetAllClassicView() []string {
	return r.AllClassicView
}

func (r *ProxyRepository) GetHTTPClassicView() []string {
	return r.HTTPClassicView
}

func (r *ProxyRepository) GetHTTPSClassicView() []string {
	return r.HTTPSClassicView
}

func (r *ProxyRepository) GetSOCKS4ClassicView() []string {
	return r.SOCKS4ClassicView
}

func (r *ProxyRepository) GetSOCKS5ClassicView() []string {
	return r.SOCKS5ClassicView
}

func (r *ProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	return r.AllAdvancedView
}

func (r *ProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	return r.HTTPAdvancedView
}

func (r *ProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	return r.HTTPSAdvancedView
}

func (r *ProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	return r.SOCKS4AdvancedView
}

func (r *ProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	return r.SOCKS5AdvancedView
}
