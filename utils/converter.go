package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/jung-kurt/gofpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func ToPDF(pathStr string) (string, error) {
	extension := strings.ToLower(filepath.Ext(pathStr))
	if extension != "" {
		extension = extension[1:] // remove dot
	}

	supportedFormats := []string{"txt", "text", "png", "jpg", "jpeg", "html", "docx", "xls", "pptx", "ppsx", "potx"}
	isSupported := false
	for _, f := range supportedFormats {
		if f == extension {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return "", fmt.Errorf("unsupported format: %s", extension)
	}

	timestamp := time.Now().Format("20060102_150405")
	dir := "dist/converter"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	outputPath := filepath.Join(dir, fmt.Sprintf("%s.pdf", timestamp))

	var err error
	switch extension {
	case "txt", "text":
		err = textToPDF(pathStr, outputPath)
	case "png", "jpg", "jpeg":
		err = imageToPDF(pathStr, outputPath)
	case "html":
		err = htmlToPDF(pathStr, outputPath)
	case "docx", "xls", "pptx", "ppsx", "potx":
		return "", fmt.Errorf("office format conversion (%s) is currently not implemented in this version due to lack of pure-go renderers", extension)
	default:
		return "", fmt.Errorf("unsupported format: %s", extension)
	}

	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func textToPDF(input, output string) error {
	content, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)
	
	y := 280.0
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if y < 10.0 {
			break
		}
		pdf.Text(10, y, line)
		y -= 15.0
	}

	return pdf.OutputFileAndClose(output)
}

func imageToPDF(input, output string) error {
	conf := model.NewDefaultConfiguration()
	return api.ImportImagesFile([]string{input}, output, nil, conf)
}

func htmlToPDF(input, output string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	absPath, err := filepath.Abs(input)
	if err != nil {
		return err
	}
	url := "file://" + absPath

	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err
	}

	return os.WriteFile(output, buf, 0644)
}
