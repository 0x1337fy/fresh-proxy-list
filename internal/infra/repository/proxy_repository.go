package repository

import (
	"cmp"
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
	Store(proxy entity.Proxy)
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
	return &ProxyRepository{}
}

func (r *ProxyRepository) Store(proxy entity.Proxy) {
	r.mu.Lock()
	defer r.mu.Unlock()

	updateProxyAll := func(proxy entity.Proxy, classicList *[]string, advancedList *[]entity.AdvancedProxy) {
		sort.Slice(*advancedList, func(i, j int) bool {
			return cmp.Compare((*advancedList)[i].Proxy, (*advancedList)[j].Proxy) < 0
		})
		n, found := slices.BinarySearchFunc(*advancedList, entity.AdvancedProxy{Proxy: proxy.Proxy}, func(a, b entity.AdvancedProxy) int {
			return cmp.Compare(a.Proxy, b.Proxy)
		})
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
		r.httpClassicView = append(r.httpClassicView, proxy.Proxy)
		r.httpAdvancedView = append(r.httpAdvancedView, proxy)
	case "HTTPS":
		r.httpsClassicView = append(r.httpsClassicView, proxy.Proxy)
		r.httpsAdvancedView = append(r.httpsAdvancedView, proxy)
	case "SOCKS4":
		r.socks4ClassicView = append(r.socks4ClassicView, proxy.Proxy)
		r.socks4AdvancedView = append(r.socks4AdvancedView, proxy)
	case "SOCKS5":
		r.socks5ClassicView = append(r.socks5ClassicView, proxy.Proxy)
		r.socks5AdvancedView = append(r.socks5AdvancedView, proxy)
	}

	updateProxyAll(proxy, &r.allClassicView, &r.allAdvancedView)
}

func (r *ProxyRepository) GetAllClassicView() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.allClassicView
}

func (r *ProxyRepository) GetHTTPClassicView() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.httpClassicView
}

func (r *ProxyRepository) GetHTTPSClassicView() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.httpsClassicView
}

func (r *ProxyRepository) GetSOCKS4ClassicView() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.socks4ClassicView
}

func (r *ProxyRepository) GetSOCKS5ClassicView() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.socks5ClassicView
}

func (r *ProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.allAdvancedView
}

func (r *ProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.httpAdvancedView
}

func (r *ProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.httpsAdvancedView
}

func (r *ProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.socks4AdvancedView
}

func (r *ProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.socks5AdvancedView
}
