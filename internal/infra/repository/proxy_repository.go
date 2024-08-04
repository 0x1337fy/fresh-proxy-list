package repository

import (
	"fresh-proxy-list/internal/entity"
	"slices"
	"sort"
	"sync"
)

type ProxyRepository struct {
	mu                 sync.RWMutex
	allClassicView     []string
	httpClassicView    []string
	httpsClassicView   []string
	socks4ClassicView  []string
	socks5ClassicView  []string
	allAdvancedView    []entity.AdvancedProxy
	httpAdvancedView   []entity.Proxy
	httpsAdvancedView  []entity.Proxy
	socks4AdvancedView []entity.Proxy
	socks5AdvancedView []entity.Proxy
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
		mu: sync.RWMutex{},
	}
}

func (r *ProxyRepository) Store(proxy *entity.Proxy) {
	r.mu.Lock()
	defer r.mu.Unlock()

	findIndex := func(slice *[]entity.AdvancedProxy, target *entity.AdvancedProxy) (int, bool) {
		for i, item := range *slice {
			if item.Proxy == target.Proxy {
				return i, true
			}
		}
		return -1, false
	}

	updateProxyAll := func(proxy *entity.Proxy, classicList *[]string, advancedList *[]entity.AdvancedProxy) {
		n, found := findIndex(advancedList, &entity.AdvancedProxy{Proxy: proxy.Proxy})
		if found {
			if proxy.Category == "HTTP" {
				(*advancedList)[n].TimeTaken = proxy.TimeTaken
			}

			if _, found := slices.BinarySearch((*advancedList)[n].Categories, proxy.Category); !found {
				(*advancedList)[n].Categories = append((*advancedList)[n].Categories, proxy.Category)
				sort.Strings((*advancedList)[n].Categories)
			}
		} else {
			*classicList = append(*classicList, proxy.Proxy)
			*advancedList = append(*advancedList, entity.AdvancedProxy{
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
			httpClassicView  = &r.httpClassicView
			httpAdvancedView = &r.httpAdvancedView
		)
		*httpClassicView = append(*httpClassicView, proxy.Proxy)
		*httpAdvancedView = append(*httpAdvancedView, *proxy)
	case "HTTPS":
		var (
			httpsClassicView  = &r.httpsClassicView
			httpsAdvancedView = &r.httpsAdvancedView
		)
		*httpsClassicView = append(*httpsClassicView, proxy.Proxy)
		*httpsAdvancedView = append(*httpsAdvancedView, *proxy)
	case "SOCKS4":
		var (
			socks4ClassicView  = &r.socks4ClassicView
			socks4AdvancedView = &r.socks4AdvancedView
		)
		*socks4ClassicView = append(*socks4ClassicView, proxy.Proxy)
		*socks4AdvancedView = append(*socks4AdvancedView, *proxy)
	case "SOCKS5":
		var (
			socks5ClassicView  = &r.socks5ClassicView
			socks5AdvancedView = &r.socks5AdvancedView
		)
		*socks5ClassicView = append(*socks5ClassicView, proxy.Proxy)
		*socks5AdvancedView = append(*socks5AdvancedView, *proxy)
	}

	updateProxyAll(proxy, &r.allClassicView, &r.allAdvancedView)
}

func (r *ProxyRepository) GetAllClassicView() []string {
	return r.allClassicView
}

func (r *ProxyRepository) GetHTTPClassicView() []string {
	return r.httpClassicView
}

func (r *ProxyRepository) GetHTTPSClassicView() []string {
	return r.httpsClassicView
}

func (r *ProxyRepository) GetSOCKS4ClassicView() []string {
	return r.socks4ClassicView
}

func (r *ProxyRepository) GetSOCKS5ClassicView() []string {
	return r.socks5ClassicView
}

func (r *ProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	return r.allAdvancedView
}

func (r *ProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	return r.httpAdvancedView
}

func (r *ProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	return r.httpsAdvancedView
}

func (r *ProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	return r.socks4AdvancedView
}

func (r *ProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	return r.socks5AdvancedView
}
