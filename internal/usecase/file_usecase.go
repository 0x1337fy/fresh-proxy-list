package usecase

import (
	"fresh-proxy-list/internal/infra/repository"
	"path/filepath"
	"strings"
	"sync"
)

type fileUsecase struct {
	fileRepository       repository.FileRepositoryInterface
	proxyRepository      repository.ProxyRepositoryInterface
	fileOutputExtensions []string
	wg                   sync.WaitGroup
}

type FileUsecaseInterface interface {
	SaveFiles()
}

func NewFileUsecase(fileRepository repository.FileRepositoryInterface, proxyRepository repository.ProxyRepositoryInterface, fileOutputExtensions []string) FileUsecaseInterface {
	return &fileUsecase{
		fileRepository:       fileRepository,
		proxyRepository:      proxyRepository,
		fileOutputExtensions: fileOutputExtensions,
		wg:                   sync.WaitGroup{},
	}
}

func (uc *fileUsecase) SaveFiles() {
	createFile := func(filename string, classic []string, advanced interface{}) {
		uc.wg.Add((len(uc.fileOutputExtensions) * 2) + 1)

		filename = strings.ToLower(filename)
		for _, ext := range uc.fileOutputExtensions {
			go func(ext string) {
				defer uc.wg.Done()
				uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+"."+ext), classic, ext)
			}(ext)
			go func(ext string) {
				defer uc.wg.Done()
				uc.fileRepository.SaveFile(filepath.Join("storage", "advanced", filename+"."+ext), advanced, ext)
			}(ext)
		}

		go func() {
			defer uc.wg.Done()
			uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+".txt"), classic, "txt")
		}()
	}

	createFile("all", uc.proxyRepository.GetAllClassicView(), uc.proxyRepository.GetAllAdvancedView())
	createFile("http", uc.proxyRepository.GetHTTPClassicView(), uc.proxyRepository.GetHTTPAdvancedView())
	createFile("https", uc.proxyRepository.GetHTTPSClassicView(), uc.proxyRepository.GetHTTPSAdvancedView())
	createFile("socks4", uc.proxyRepository.GetSOCKS4ClassicView(), uc.proxyRepository.GetSOCKS4AdvancedView())
	createFile("socks5", uc.proxyRepository.GetSOCKS5ClassicView(), uc.proxyRepository.GetSOCKS5AdvancedView())
	uc.wg.Wait()
}
