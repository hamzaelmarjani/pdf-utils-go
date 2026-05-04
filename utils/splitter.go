package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Split(path string) ([]string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dir := filepath.Join("dist/spliter", timestamp)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	conf := model.NewDefaultConfiguration()
	if err := api.SplitFile(path, dir, 1, conf); err != nil {
		return nil, err
	}

	// Read the directory and rename files to match Rust version (1.pdf, 2.pdf, ...)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var outputPaths []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pdf" {
			// pdfcpu naming: <filename>_<page>.pdf
			// We want: <page>.pdf
			oldPath := filepath.Join(dir, file.Name())
			
			// Find the last underscore
			name := strings.TrimSuffix(file.Name(), ".pdf")
			lastUnderscore := strings.LastIndex(name, "_")
			if lastUnderscore != -1 {
				pageNum := name[lastUnderscore+1:]
				newPath := filepath.Join(dir, pageNum+".pdf")
				if err := os.Rename(oldPath, newPath); err == nil {
					outputPaths = append(outputPaths, newPath)
				} else {
					outputPaths = append(outputPaths, oldPath)
				}
			} else {
				outputPaths = append(outputPaths, oldPath)
			}
		}
	}

	return outputPaths, nil
}
