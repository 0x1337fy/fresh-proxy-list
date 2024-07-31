package usecase

import (
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/internal/service"
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

	ipPort := strings.Split(proxy, ":")
	if len(ipPort) != 2 {
		return fmt.Errorf("proxy format incorrect")
	}

	_, loaded := uc.proxyMap.LoadOrStore(source.Category+"_"+proxy, true)
	if loaded {
		return fmt.Errorf("proxy has been processed")
	}

	var (
		data *entity.Proxy
		err  error
	)
	ip, port := ipPort[0], ipPort[1]
	if source.IsChecked {
		data, err = uc.proxyService.Check(source.Category, ip, port)
		if err != nil {
			return err
		}
	} else {
		data = &entity.Proxy{
			Proxy:     proxy,
			IP:        ip,
			Port:      port,
			Category:  source.Category,
			TimeTaken: 0,
			CheckedAt: "",
		}
	}
	uc.proxyRepository.Store(*data)

	return nil
}

func (uc *ProxyUsecase) GetAllAdvancedView() []entity.AdvancedProxy {
	return uc.proxyRepository.GetAllAdvancedView()
}
