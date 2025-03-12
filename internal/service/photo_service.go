package service

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"

	"wealthlist/internal/logger"
	"wealthlist/internal/repo"
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

func (s *PhotoService) UploadPhoto(millionaireID int, file *multipart.FileHeader) (string, error) {
	randSuffix := strconv.Itoa(rand.Intn(10000))
	uniqueFileName := fmt.Sprintf("%d_%s_%s", millionaireID, randSuffix, path.Base(file.Filename))
	savePath := "uploads/photos/" + uniqueFileName

	if err := os.MkdirAll("uploads/photos", os.ModePerm); err != nil {
		s.log.Error("Error creating directory", logger.Err(err))
		return "", err
	}

	if err := saveUploadedFile(file, savePath); err != nil {
		s.log.Error("Error saving file", logger.Err(err))
		return "", err
	}

	return savePath, nil
}

func (s *PhotoService) UpdatePhoto(millionaireID int, filePath string) error {
	return s.photoRepo.UpdatePhotoPath(millionaireID, filePath)
}

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

func (s *PhotoService) GetPhotoPath(millionaireID int) (string, error) {
	return s.photoRepo.GetPhotoPath(millionaireID)
}

func (s *PhotoService) ClearPhotoPath(millionaireID int) error {
	return s.photoRepo.ClearPhotoPath(millionaireID)
}
