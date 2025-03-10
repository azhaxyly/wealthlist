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
	const op = "service.CreateMillionaire"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating millionaire")

	err := s.repo.Create(m)
	if err != nil {
		log.Error("failed to create millionaire", logger.Err(err))
		return err
	}

	log.Info("millionaire created successfully", slog.Int("id", m.ID))
	return nil
}

func (s *millionaireService) SearchMillionaire(lastName, firstName, middleName, country string, pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	const op = "service.SearchMillionaire"
	log := s.log.With(
		slog.String("op", op),
		slog.String("lastName", lastName),
		slog.String("firstName", firstName),
		slog.Int("pageNum", pageNum),
		slog.Int("pageSize", pageSize),
	)

	log.Debug("starting search")

	if pageNum < 1 {
		pageNum = 1
		log.Warn("pageNum adjusted to minimum value", slog.Int("newPageNum", pageNum))
	}
	if pageSize < 1 {
		pageSize = 10
		log.Warn("pageSize adjusted to default", slog.Int("newPageSize", pageSize))
	}

	filter := repo.MillionaireFilter{
		LastName:   lastName,
		FirstName:  firstName,
		MiddleName: middleName,
		Country:    country,
	}

	result, err := s.repo.Search(filter, pageNum, pageSize)
	if err != nil {
		log.Error("search failed", logger.Err(err))
		return models.PaginationMillionaireDto{}, err
	}

	log.Info("search completed",
		slog.Int("totalResults", result.Total),
		slog.Int("returnedResults", len(result.Millionaires)),
	)
	return result, nil
}

func (s *millionaireService) GetAllMillionaires(pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	const op = "service.GetAllMillionaires"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("pageNum", pageNum),
		slog.Int("pageSize", pageSize),
	)

	log.Debug("fetching millionaires")

	if pageNum < 1 {
		pageNum = 1
		log.Warn("pageNum adjusted to minimum value", slog.Int("newPageNum", pageNum))
	}
	if pageSize < 1 {
		pageSize = 10
		log.Warn("pageSize adjusted to default", slog.Int("newPageSize", pageSize))
	}

	result, err := s.repo.GetAll(pageNum, pageSize)
	if err != nil {
		log.Error("failed to fetch millionaires", logger.Err(err))
		return models.PaginationMillionaireDto{}, err
	}

	log.Info("successfully fetched millionaires",
		slog.Int("total", result.Total),
		slog.Int("returned", len(result.Millionaires)),
	)
	return result, nil
}

func (s *millionaireService) GetMillionaireByID(id int) (*models.Millionaire, error) {
	const op = "service.GetMillionaireByID"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("id", id),
	)

	log.Debug("fetching millionaire")

	millionaire, err := s.repo.GetByID(id)
	if err != nil {
		log.Error("failed to fetch millionaire", logger.Err(err))
		return nil, err
	}

	if millionaire == nil {
		log.Warn("millionaire not found")
		return nil, nil
	}

	log.Debug("successfully fetched millionaire")
	return millionaire, nil
}

// func (s *MillionaireService) GetHomePageData() (models.HomePageDto, error) {
// 	millionaires, err := s.repo.GetWithPhotos()
// 	if err != nil {
// 		return models.HomePageDto{}, err
// 	}

// 	topMillionaires := make([]models.Millionaire, 0)
// 	featured := make([]models.Millionaire, 0)

// 	for i, m := range millionaires {
// 		if i < 5 {
// 			topMillionaires = append(topMillionaires, m)
// 		} else {
// 			featured = append(featured, m)
// 		}
// 	}

// 	return models.HomePageDto{
// 		TopMillionaires: topMillionaires,
// 		Featured:        featured,
// 	}, nil
// }

func (s *millionaireService) UpdateMillionaire(m *models.Millionaire) error {
	const op = "service.UpdateMillionaire"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("id", m.ID),
	)

	log.Info("updating millionaire")

	err := s.repo.Update(m)
	if err != nil {
		log.Error("update failed", logger.Err(err))
		return err
	}

	log.Info("millionaire updated successfully")
	return nil
}

func (s *millionaireService) DeleteMillionaire(id int) error {
	const op = "service.DeleteMillionaire"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("id", id),
	)

	log.Info("deleting millionaire")

	err := s.repo.Delete(id)
	if err != nil {
		log.Error("delete failed", logger.Err(err))
		return err
	}

	log.Info("millionaire deleted successfully")
	return nil
}
