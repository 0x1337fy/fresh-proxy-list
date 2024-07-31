package usecase

import (
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/config"
	"fresh-proxy-list/internal/infra/repository"
	"path/filepath"
	"strings"
)

type fileUsecase struct {
	fileRepository       repository.FileRepositoryInterface
	proxyRepository      repository.ProxyRepositoryInterface
	fileOutputExtensions []string
}

type FileUsecaseInterface interface {
	SaveFiles()
}

func NewFileUsecase(fileRepository repository.FileRepositoryInterface, proxyRepository repository.ProxyRepositoryInterface) FileUsecaseInterface {
	return &fileUsecase{
		fileRepository:       fileRepository,
		proxyRepository:      proxyRepository,
		fileOutputExtensions: config.FileOutputExtensions,
	}
}

func (uc *fileUsecase) SaveFiles() {
	createFileForAllCategories := func(filename string, classic []string, advanced []entity.AdvancedProxy) {
		filename = strings.ToLower(filename)
		for _, ext := range uc.fileOutputExtensions {
			uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+"."+ext), classic, ext)
			uc.fileRepository.SaveFile(filepath.Join("storage", "advanced", filename+"."+ext), advanced, ext)
		}
		uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+".txt"), classic, "txt")
	}

	createFile := func(filename string, classic []string, advanced []entity.Proxy) {
		filename = strings.ToLower(filename)
		for _, ext := range uc.fileOutputExtensions {
			uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+"."+ext), classic, ext)
			uc.fileRepository.SaveFile(filepath.Join("storage", "advanced", filename+"."+ext), advanced, ext)
		}
		uc.fileRepository.SaveFile(filepath.Join("storage", "classic", filename+".txt"), classic, "txt")
	}

	createFileForAllCategories("all", uc.proxyRepository.GetAllClassicView(), uc.proxyRepository.GetAllAdvancedView())
	createFile("http", uc.proxyRepository.GetHTTPClassicView(), uc.proxyRepository.GetHTTPAdvancedView())
	createFile("https", uc.proxyRepository.GetHTTPSClassicView(), uc.proxyRepository.GetHTTPSAdvancedView())
	createFile("socks4", uc.proxyRepository.GetSOCKS4ClassicView(), uc.proxyRepository.GetSOCKS4AdvancedView())
	createFile("socks5", uc.proxyRepository.GetSOCKS5ClassicView(), uc.proxyRepository.GetSOCKS5AdvancedView())
}
