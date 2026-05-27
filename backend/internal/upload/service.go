package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	baseDir string
	baseURL string
}

type SaveImageResult struct {
	URL      string
	SavePath string
}

func NewService() *Service {
	return &Service{
		baseDir: "uploads/products",
		baseURL: "/uploads/products",
	}
}

func (s *Service) PrepareImage(file *multipart.FileHeader) (*SaveImageResult, error) {
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return nil, fmt.Errorf("only jpg, jpeg, png, webp are allowed")
	}

	if file.Size > 5*1024*1024 {
		return nil, fmt.Errorf("file size must be less than 5MB")
	}

	if err := os.MkdirAll(s.baseDir, os.ModePerm); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(s.baseDir, filename)

	return &SaveImageResult{
		URL:      s.baseURL + "/" + filename,
		SavePath: savePath,
	}, nil
}
