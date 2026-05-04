package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ledongthuc/pdf"
)

func FromPDF(pathStr string, format string) (string, error) {
	format = strings.ToLower(format)
	supportedFormats := []string{"text", "png", "jpg", "jpeg", "html"}
	isSupported := false
	for _, f := range supportedFormats {
		if f == format {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	timestamp := time.Now().Format("20060102_150405")
	dir := "dist/reverser"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	outputPath := filepath.Join(dir, fmt.Sprintf("%s.%s", timestamp, format))

	var err error
	switch format {
	case "text":
		err = pdfToText(pathStr, outputPath)
	case "html":
		err = pdfToHTML(pathStr, outputPath)
	case "png", "jpg", "jpeg":
		err = pdfToImage(pathStr, outputPath, format)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func pdfToText(input, output string) error {
	content, err := getPDFText(input)
	if err != nil {
		return err
	}
	return os.WriteFile(output, []byte(content), 0644)
}

func pdfToHTML(input, output string) error {
	content, err := getPDFText(input)
	if err != nil {
		return err
	}

	htmlContent := fmt.Sprintf("<html><body><pre>%s</pre></body></html>", escapeHTML(content))
	return os.WriteFile(output, []byte(htmlContent), 0644)
}

func getPDFText(input string) (string, error) {
	f, r, err := pdf.Open(input)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var text strings.Builder
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			var lastX float64
			for i, word := range row.Content {
				if i > 0 && word.X > lastX+word.W {
					text.WriteString(" ")
				}
				text.WriteString(word.S)
				lastX = word.X + word.W
			}
			text.WriteString("\n")
		}
	}
	return text.String(), nil
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func pdfToImage(input, output, format string) error {
	// Try Chrome first
	err := pdfToImageWithChrome(input, output, format)
	if err == nil {
		return nil
	}

	// Fallback to native tools
	return pdfToImageWithNativeTool(input, output, format)
}

func pdfToImageWithChrome(input, output, format string) error {
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
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return err
	}

	return os.WriteFile(output, buf, 0644)
}

func pdfToImageWithNativeTool(input, output, format string) error {
	sipsFormat := "jpeg"
	if format == "png" {
		sipsFormat = "png"
	}

	// Try sips (macOS)
	cmd := exec.Command("sips", "-s", "format", sipsFormat, input, "--out", output)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try magick (ImageMagick 7)
	firstPage := input + "[0]"
	cmd = exec.Command("magick", firstPage, output)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try convert (ImageMagick 6)
	cmd = exec.Command("convert", firstPage, output)
	if err := cmd.Run(); err == nil {
		return nil
	}

	return fmt.Errorf("no available PDF renderer found. Install Chrome/Chromium, macOS sips, or ImageMagick")
}
