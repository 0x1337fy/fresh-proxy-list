package main

import (
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/config"
	"fresh-proxy-list/internal/infra/repository"
	"fresh-proxy-list/internal/usecase"
	"fresh-proxy-list/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Runners struct {
	fetcherUtil      utils.FetcherUtilInterface
	urlParserUtil    utils.URLParserUtilInterface
	sourceRepository repository.SourceRepositoryInterface
	proxyRepository  repository.ProxyRepositoryInterface
	fileRepository   repository.FileRepositoryInterface
}

func main() {
	if err := runApplication(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func runApplication() error {
	runtime.GOMAXPROCS(2)
	loadEnv()

	mkdirAll := func(path string, perm os.FileMode) error {
		return os.MkdirAll(path, perm)
	}
	create := func(name string) (io.Writer, error) {
		file, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	client := http.DefaultTransport
	fetcherUtil := utils.NewFetcher(client, createHTTPRequest)
	urlParserUtil := utils.NewURLParser()
	sourceRepository := repository.NewSourceRepository(os.Getenv("PROXY_RESOURCES"))
	proxyRepository := repository.NewProxyRepository()
	fileRepository := repository.NewFileRepository(mkdirAll, create)

	runners := Runners{
		fetcherUtil:      fetcherUtil,
		urlParserUtil:    urlParserUtil,
		sourceRepository: sourceRepository,
		proxyRepository:  proxyRepository,
		fileRepository:   fileRepository,
	}

	return run(runners)
}

func loadEnv() error {
	return godotenv.Load()
}

func createHTTPRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func run(runners Runners) error {
	start := time.Now()

	sourceUsecase := usecase.NewSourceUsecase(runners.sourceRepository)
	sources, err := sourceUsecase.LoadSources()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	proxyCategories := config.ProxyCategories
	proxyUsecase := usecase.NewProxyUsecase(runners.proxyRepository, runners.fetcherUtil, runners.urlParserUtil)
	for i, source := range sources {
		if _, found := slices.BinarySearch(proxyCategories, source.Category); found {
			wg.Add(1)
			go func(source entity.Source) {
				defer wg.Done()

				body, err := runners.fetcherUtil.FetchData(source.URL)
				if err != nil {
					return
				}

				var (
					innerWG sync.WaitGroup
					proxies []string
				)
				switch source.Method {
				case "LIST":
					proxies = strings.Split(strings.TrimSpace(string(body)), "\n")
				case "SCRAP":
					re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+){3}:[0-9]+`)
					proxies = re.FindAllString(string(body), -1)
				default:
					return
				}

				for _, proxy := range proxies {
					innerWG.Add(1)
					go func(source entity.Source, proxy string) {
						defer innerWG.Done()
						_ = proxyUsecase.ProcessProxy(source, proxy)
					}(source, proxy)
				}
				innerWG.Wait()
			}(source)
		} else {
			log.Printf("Index %v: proxy category not found", i)
		}
	}
	wg.Wait()

	fileUsecase := usecase.NewFileUsecase(runners.fileRepository, runners.proxyRepository)
	fileUsecase.SaveFiles()

	log.Printf("Number of proxies     : %v", len(runners.proxyRepository.GetAllAdvancedView()))
	log.Printf("Time-consuming process: %v", time.Since(start))
	return nil
}
