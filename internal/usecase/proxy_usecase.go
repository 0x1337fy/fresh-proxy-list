package usecase

import (
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/config"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/pkg/utils"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"h12.io/socks"
)

type ProxyUsecase struct {
	proxyRepository   repository.ProxyRepositoryInterface
	fetcherUtil       utils.FetcherUtilInterface
	urlParserUtil     utils.URLParserUtilInterface
	httpTestingSites  []string
	httpsTestingSites []string
	userAgents        []string
	proxyMap          sync.Map
	semaphore         chan struct{}
}

type ProxyUsecaseInterface interface {
	ProcessProxy(source entity.Source, proxy string) error
	IsProxyWorking(source entity.Source, ip string, port string) (entity.Proxy, error)
	GetTestingSite(category string) string
	GetRandomUserAgent() string
}

func NewProxyUsecase(proxyRepository repository.ProxyRepositoryInterface, fetcherUtil utils.FetcherUtilInterface, urlParserUtil utils.URLParserUtilInterface) ProxyUsecaseInterface {
	return &ProxyUsecase{
		proxyRepository:   proxyRepository,
		fetcherUtil:       fetcherUtil,
		urlParserUtil:     urlParserUtil,
		httpTestingSites:  config.HTTPTestingSites,
		httpsTestingSites: config.HTTPSTestingSites,
		userAgents:        config.UserAgents,
		proxyMap:          sync.Map{},
		semaphore:         make(chan struct{}, 500),
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
		data entity.Proxy
		err  error
	)
	ip, port := ipPort[0], ipPort[1]
	if source.IsChecked {
		uc.semaphore <- struct{}{}
		defer func() { <-uc.semaphore }()

		data, err = uc.IsProxyWorking(source, ip, port)
		if err != nil {
			return err
		}
	} else {
		data = entity.Proxy{
			Proxy:     proxy,
			IP:        ip,
			Port:      port,
			Category:  source.Category,
			TimeTaken: 0,
			CheckedAt: "",
		}
	}
	uc.proxyRepository.Store(data)

	return nil
}

func (uc *ProxyUsecase) IsProxyWorking(source entity.Source, ip string, port string) (entity.Proxy, error) {
	var (
		transport   *http.Transport
		proxy       = ip + ":" + port
		proxyURI    = strings.ToLower(source.Category + "://" + proxy)
		testingSite = uc.GetTestingSite(source.Category)
		timeout     = 5 * time.Second
	)

	if source.Category == "HTTP" || source.Category == "HTTPS" {
		proxyURL, err := uc.urlParserUtil.Parse(proxyURI)
		if err != nil {
			return entity.Proxy{}, fmt.Errorf("error parsing proxy URL: %v", err)
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
		}
		if source.Category == "HTTPS" {
			transport.TLSHandshakeTimeout = timeout
		}
	} else if source.Category == "SOCKS4" || source.Category == "SOCKS5" {
		proxyURL := socks.Dial(proxyURI)
		transport = &http.Transport{
			Dial: proxyURL,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
		}
	} else {
		return entity.Proxy{}, fmt.Errorf("proxy category %s not supported", source.Category)
	}

	uc.fetcherUtil.SetClient(transport)
	req, err := uc.fetcherUtil.NewRequest("GET", testingSite, nil)
	if err != nil {
		return entity.Proxy{}, err
	}
	req.Header.Set("User-Agent", uc.GetRandomUserAgent())

	startTime := time.Now()
	resp, err := uc.fetcherUtil.Do(req)
	if err != nil {
		return entity.Proxy{}, fmt.Errorf("request error: %s", err)
	}
	defer resp.Body.Close()
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Seconds()

	if resp.StatusCode != http.StatusOK {
		return entity.Proxy{}, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return entity.Proxy{
		Proxy:     proxy,
		IP:        ip,
		Port:      port,
		Category:  source.Category,
		CheckedAt: endTime.Format(time.RFC3339),
		TimeTaken: timeTaken,
	}, nil
}

func (uc *ProxyUsecase) GetTestingSite(category string) string {
	if category == "HTTPS" {
		return uc.httpsTestingSites[rand.Intn(len(uc.httpsTestingSites))]
	}
	return uc.httpTestingSites[rand.Intn(len(uc.httpTestingSites))]
}

func (uc *ProxyUsecase) GetRandomUserAgent() string {
	return uc.userAgents[rand.Intn(len(uc.userAgents))]
}

func (uc *ProxyUsecase) GetAllAdvancedView() []entity.AdvancedProxy {
	return uc.proxyRepository.GetAllAdvancedView()
}
