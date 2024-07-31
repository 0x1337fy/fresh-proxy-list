package usecase

import (
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/internal/infra/repository"
)

type sourceUsecase struct {
	sourceRepository repository.SourceRepositoryInterface
}

type SourceUsecaseInterface interface {
	LoadSources() ([]entity.Source, error)
}

func NewSourceUsecase(sourceRepository repository.SourceRepositoryInterface) SourceUsecaseInterface {
	return &sourceUsecase{
		sourceRepository: sourceRepository,
	}
}

func (uc *sourceUsecase) LoadSources() ([]entity.Source, error) {
	return uc.sourceRepository.LoadSources()
}
