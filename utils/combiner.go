package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Combine(paths []string, output string) (string, error) {
	if len(paths) < 2 {
		return "", fmt.Errorf("at least 2 PDF files are required to combine")
	}

	outputPath := output
	if outputPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		dir := "dist/combiner"
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
		outputPath = filepath.Join(dir, fmt.Sprintf("%s.pdf", timestamp))
	}

	conf := model.NewDefaultConfiguration()
	if err := api.MergeCreateFile(paths, outputPath, false, conf); err != nil {
		return "", err
	}

	return outputPath, nil
}
