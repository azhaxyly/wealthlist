package service

import (
	"log/slog"
	"wealthlist/internal/logger"
	"wealthlist/internal/models"
	"wealthlist/internal/repo"
)

type HomeService struct {
	repo repo.MillionaireRepository
	log  *slog.Logger
}

func NewHomeService(repo repo.MillionaireRepository, log *slog.Logger) *HomeService {
	return &HomeService{repo: repo, log: log}
}

func (s *HomeService) GetHomePageData(baseURL string) (*models.HomePageDto, error) {
	topMillionaires, err := s.repo.GetTopMillionaires(baseURL)
	if err != nil {
		s.log.Error("Error fetching top millionaires", logger.Err(err))
		return nil, err
	}

	return &models.HomePageDto{
		TopMillionaires: topMillionaires,
	}, nil
}
