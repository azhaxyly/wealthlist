package service

import (
	"millionaire-list/internal/models"
	"millionaire-list/internal/repo"
)

type MillionaireService struct {
	repo repo.MillionaireRepository
}

func NewMillionaireService(repo repo.MillionaireRepository) *MillionaireService {
	return &MillionaireService{repo: repo}
}

func (s *MillionaireService) CreateMillionaire(m *models.Millionaire) error {
	return s.repo.Create(m)
}

func (s *MillionaireService) SearchMillionaire(lastName, firstName, middleName, country string, pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filter := repo.MillionaireFilter{
		LastName:   lastName,
		FirstName:  firstName,
		MiddleName: middleName,
		Country:    country,
	}

	return s.repo.Search(filter, pageNum, pageSize)
}

func (s *MillionaireService) GetAllMillionaires(pageNum, pageSize int) (models.PaginationMillionaireDto, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return s.repo.GetAll(pageNum, pageSize)
}

func (s *MillionaireService) GetMillionaireByID(id int) (*models.Millionaire, error) {
	return s.repo.GetByID(id)
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

func (s *MillionaireService) UpdateMillionaire(m *models.Millionaire) error {
	return s.repo.Update(m)
}

func (s *MillionaireService) DeleteMillionaire(id int) error {
	return s.repo.Delete(id)
}
