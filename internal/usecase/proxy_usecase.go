package usecase

import (
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/internal/service"
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type ProxyUsecase struct {
	proxyRepository repository.ProxyRepositoryInterface
	proxyService    service.ProxyServiceInterface
	proxyMap        sync.Map
}

type ProxyUsecaseInterface interface {
	ProcessProxy(source entity.Source, proxy string) error
	IsSpecialIP(ip string) bool
	ParseCIDR(cidr string) *net.IPNet
	GetAllAdvancedView() []entity.AdvancedProxy
}

func NewProxyUsecase(proxyRepository repository.ProxyRepositoryInterface, proxyService service.ProxyServiceInterface) ProxyUsecaseInterface {
	return &ProxyUsecase{
		proxyRepository: proxyRepository,
		proxyService:    proxyService,
		proxyMap:        sync.Map{},
	}
}

func (uc *ProxyUsecase) ProcessProxy(source entity.Source, proxy string) error {
	proxy = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(proxy, "\r", ""), "\n", ""))
	if proxy == "" {
		return fmt.Errorf("proxy not found")
	}

	proxyParts := strings.Split(proxy, ":")
	if len(proxyParts) != 2 {
		return fmt.Errorf("proxy format incorrect")
	}

	pattern := `^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\:(0|[1-9][0-9]{0,4})$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(proxy) {
		return fmt.Errorf("proxy format not match")
	}

	if uc.IsSpecialIP(proxyParts[0]) {
		return fmt.Errorf("proxy belongs to special ip")
	}

	port, err := strconv.Atoi(proxyParts[1])
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("proxy port format incorrect")
	}

	_, loaded := uc.proxyMap.LoadOrStore(source.Category+"_"+proxy, true)
	if loaded {
		return fmt.Errorf("proxy has been processed")
	}

	var data *entity.Proxy
	proxyIP, proxyPort := proxyParts[0], proxyParts[1]
	if source.IsChecked {
		data, err = uc.proxyService.Check(source.Category, proxyIP, proxyPort)
		if err != nil {
			return err
		}
	} else {
		data = &entity.Proxy{
			Proxy:     proxy,
			IP:        proxyIP,
			Port:      proxyPort,
			Category:  source.Category,
			TimeTaken: 0,
			CheckedAt: "",
		}
	}
	uc.proxyRepository.Store(*data)

	return nil
}

func (uc *ProxyUsecase) IsSpecialIP(ip string) bool {
	if _, found := slices.BinarySearch([]string{
		"0.0.0.0",
		"127.0.0.1",
		"255.255.255.255",
	}, ip); found {
		return true
	}

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return true
	}

	privateRanges := []struct {
		netIP *net.IPNet
	}{
		{netIP: uc.ParseCIDR("10.0.0.0/8")},
		{netIP: uc.ParseCIDR("172.16.0.0/12")},
		{netIP: uc.ParseCIDR("192.168.0.0/16")},
		{netIP: uc.ParseCIDR("169.254.0.0/16")}, // link-local
		{netIP: uc.ParseCIDR("240.0.0.0/4")},    // reserved for special use
		{netIP: uc.ParseCIDR("224.0.0.0/4")},    // multicast
	}
	for _, r := range privateRanges {
		if r.netIP.Contains(ipAddr) {
			return true
		}
	}

	return false
}

func (uc *ProxyUsecase) ParseCIDR(cidr string) *net.IPNet {
	_, netIP, _ := net.ParseCIDR(cidr)
	return netIP
}

func (uc *ProxyUsecase) GetAllAdvancedView() []entity.AdvancedProxy {
	return uc.proxyRepository.GetAllAdvancedView()
}
