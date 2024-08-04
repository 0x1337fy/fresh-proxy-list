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
	specialIPs      []string
	privateIPs      []net.IPNet
}

type ProxyUsecaseInterface interface {
	ProcessProxy(source entity.Source, proxy string) error
	IsSpecialIP(ip string) bool
	GetAllAdvancedView() []entity.AdvancedProxy
}

func NewProxyUsecase(
	proxyRepository repository.ProxyRepositoryInterface,
	proxyService service.ProxyServiceInterface,
	specialIPs []string,
	privateIPs []net.IPNet,
) ProxyUsecaseInterface {
	return &ProxyUsecase{
		proxyRepository: proxyRepository,
		proxyService:    proxyService,
		specialIPs:      specialIPs,
		privateIPs:      privateIPs,
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
	if _, found := slices.BinarySearch(uc.specialIPs, ip); found {
		return true
	}

	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return true
	}

	if ipAddress.IsLoopback() || ipAddress.IsMulticast() || ipAddress.IsUnspecified() {
		return true
	}

	for _, r := range uc.privateIPs {
		if r.Contains(ipAddress) {
			return true
		}
	}

	return false
}

func (uc *ProxyUsecase) GetAllAdvancedView() []entity.AdvancedProxy {
	return uc.proxyRepository.GetAllAdvancedView()
}
