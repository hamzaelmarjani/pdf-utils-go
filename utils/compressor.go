package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Compress(path string, percentage float32) (string, error) {
	if percentage < 0 || percentage > 1 {
		return "", fmt.Errorf("percentage must be between 0 and 1")
	}

	timestamp := time.Now().Format("20060102_150405")
	dir := "dist/compressor"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	outputPath := filepath.Join(dir, fmt.Sprintf("%s.pdf", timestamp))

	conf := model.NewDefaultConfiguration()
	// pdfcpu optimization doesn't directly use a percentage for image scaling in its high-level API
	// but OptimizeFile is very effective at reducing size by removing redundant objects and compressing streams.
	if err := api.OptimizeFile(path, outputPath, conf); err != nil {
		return "", err
	}

	return outputPath, nil
}
