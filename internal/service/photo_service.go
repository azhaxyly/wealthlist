package service

import (
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"

	"millionaire-list/internal/logger"
	"millionaire-list/internal/repo"

	"github.com/google/uuid"
)

type PhotoService struct {
	photoRepo *repo.PhotoRepo
	log       *slog.Logger
}

func NewPhotoService(photoRepo *repo.PhotoRepo, log *slog.Logger) *PhotoService {
	return &PhotoService{
		photoRepo: photoRepo,
		log:       log,
	}
}

// Загружает фото и возвращает путь к файлу
func (s *PhotoService) UploadPhoto(millionaireID int, file *multipart.FileHeader) (string, error) {
	uniqueFileName := uuid.New().String() + filepath.Ext(file.Filename)
	savePath := "uploads/photos/" + uniqueFileName

	// Создаём директорию, если её нет
	if err := os.MkdirAll("uploads/photos", os.ModePerm); err != nil {
		s.log.Error("Error creating directory", logger.Err(err))
		return "", err
	}

	// Сохраняем файл
	if err := saveUploadedFile(file, savePath); err != nil {
		s.log.Error("Error saving file", logger.Err(err))
		return "", err
	}

	return savePath, nil
}

// Обновляет путь к фото в базе данных
func (s *PhotoService) UpdatePhoto(millionaireID int, filePath string) error {
	return s.photoRepo.UpdatePhotoPath(millionaireID, filePath)
}

// Вспомогательная функция для сохранения файла
func saveUploadedFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
