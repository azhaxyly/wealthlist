package service

import (
	"log/slog"
	"millionaire-list/internal/logger"
	"millionaire-list/internal/models"
	"millionaire-list/internal/repo"
)

type MillionaireServiceInterface interface {
	CreateMillionaire(m *models.Millionaire) error
	SearchMillionaire(lastName, firstName, middleName, country string, pageNum, pageSize int) (models.PaginationMillionaireDto, error)
	GetAllMillionaires(pageNum, pageSize int) (models.PaginationMillionaireDto, error)
	GetMillionaireByID(id int) (*models.Millionaire, error)
	UpdateMillionaire(m *models.Millionaire) error
	DeleteMillionaire(id int) error
}

type millionaireService struct {
	repo repo.MillionaireRepository
	log  *slog.Logger
}

var _ MillionaireServiceInterface = (*millionaireService)(nil) // compile-time check

func NewMillionaireService(repo repo.MillionaireRepository, log *slog.Logger) *millionaireService {
	return &millionaireService{
		repo: repo,
		log:  log,
	}
}

func (s *millionaireService) CreateMillionaire(m *models.Millionaire) error {
	s.log.Info("Creating millionaire")

	err := s.repo.Create(m)
	if err != nil {
		s.log.Error("Failed to create millionaire", logger.Err(err))
		return err
	}

	s.log.Info("Millionaire created successfully", slog.Int("id", m.ID))
	return nil
}

func (s *millionaireService) SearchMillionaire(lastName, firstName, middleName, country string, pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	s.log.Debug("Searching millionaire",
		slog.String("lastName", lastName),
		slog.String("firstName", firstName),
		slog.Int("pageNum", pageNum),
		slog.Int("pageSize", pageSize),
	)

	if pageNum < 1 {
		pageNum = 1
		s.log.Warn("pageNum adjusted", slog.Int("newPageNum", pageNum))
	}
	if pageSize < 1 {
		pageSize = 10
		s.log.Warn("pageSize adjusted", slog.Int("newPageSize", pageSize))
	}

	filter := repo.MillionaireFilter{
		LastName:   lastName,
		FirstName:  firstName,
		MiddleName: middleName,
		Country:    country,
	}

	result, err := s.repo.Search(filter, pageNum, pageSize)
	if err != nil {
		s.log.Error("Search failed", logger.Err(err))
		return models.PaginationMillionaireDto{}, err
	}

	s.log.Info("Search completed",
		slog.Int("totalResults", result.Total),
		slog.Int("returnedResults", len(result.Millionaires)),
	)
	return result, nil
}

func (s *millionaireService) GetAllMillionaires(pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	s.log.Debug("Fetching millionaires",
		slog.Int("pageNum", pageNum),
		slog.Int("pageSize", pageSize),
	)

	if pageNum < 1 {
		pageNum = 1
		s.log.Warn("pageNum adjusted", slog.Int("newPageNum", pageNum))
	}
	if pageSize < 1 {
		pageSize = 10
		s.log.Warn("pageSize adjusted", slog.Int("newPageSize", pageSize))
	}

	result, err := s.repo.GetAll(pageNum, pageSize)
	if err != nil {
		s.log.Error("Failed to fetch millionaires", logger.Err(err))
		return models.PaginationMillionaireDto{}, err
	}

	s.log.Info("Successfully fetched millionaires",
		slog.Int("total", result.Total),
		slog.Int("returned", len(result.Millionaires)),
	)
	return result, nil
}

func (s *millionaireService) GetMillionaireByID(id int) (*models.Millionaire, error) {
	s.log.Debug("Fetching millionaire", slog.Int("id", id))

	millionaire, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Error("Failed to fetch millionaire", logger.Err(err))
		return nil, err
	}

	if millionaire == nil {
		s.log.Warn("Millionaire not found")
		return nil, nil
	}

	s.log.Debug("Successfully fetched millionaire")
	return millionaire, nil
}

func (s *millionaireService) UpdateMillionaire(m *models.Millionaire) error {
	s.log.Info("Updating millionaire", slog.Int("id", m.ID))

	err := s.repo.Update(m)
	if err != nil {
		s.log.Error("Update failed", logger.Err(err))
		return err
	}

	s.log.Info("Millionaire updated successfully")
	return nil
}

func (s *millionaireService) DeleteMillionaire(id int) error {
	s.log.Info("Deleting millionaire", slog.Int("id", id))

	err := s.repo.Delete(id)
	if err != nil {
		s.log.Error("Delete failed", logger.Err(err))
		return err
	}

	s.log.Info("Millionaire deleted successfully")
	return nil
}
