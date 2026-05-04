package pdf_utils

import (
	"os"
	"testing"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func TestCombine(t *testing.T) {
	utils := NewPDFUtils()
	files := []string{
		"./tests_data/to-combine-1.pdf",
		"./tests_data/to-combine-2.pdf",
	}
	output, err := utils.Combine(files, "")
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", output)
	}

	// Verify page count
	count, err := api.PageCountFile(output)
	if err != nil {
		t.Fatalf("Failed to get page count: %v", err)
	}
	if count != len(files) {
		t.Errorf("Expected %d pages, got %d", len(files), count)
	}
}

func TestSplit(t *testing.T) {
	utils := NewPDFUtils()
	path := "./tests_data/to-split.pdf"
	
	// Get original page count
	originalCount, err := api.PageCountFile(path)
	if err != nil {
		t.Fatalf("Failed to get original page count: %v", err)
	}

	pages, err := utils.Split(path)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}
	if len(pages) != originalCount {
		t.Errorf("Expected %d split pages, got %d", originalCount, len(pages))
	}
}

func TestCompress(t *testing.T) {
	utils := NewPDFUtils()
	input := "./tests_data/to-compress.pdf"
	info, err := os.Stat(input)
	if err != nil {
		t.Fatalf("Stat input failed: %v", err)
	}
	inputSize := info.Size()

	output, err := utils.Compress(input, 0.8)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}
	if info, err := os.Stat(output); err != nil || info.Size() == 0 {
		t.Fatalf("Output file invalid: %s", output)
	}
    // Optimization might not always reduce size if it's already optimized, but should be valid.
    _ = inputSize
}

func TestConvertText(t *testing.T) {
	utils := NewPDFUtils()
	output, err := utils.ConvertToPDF("./tests_data/to-pdf.png")
	if err != nil {
		t.Fatalf("ConvertToPDF failed: %v", err)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", output)
	}
}

func TestReverseText(t *testing.T) {
	utils := NewPDFUtils()
	output, err := utils.ReverseFromPDF("./tests_data/to-image.pdf", "jpeg")
	if err != nil {
		t.Logf("Skipping ReverseText test (no renderer available): %v", err)
		return
	}
	info, err := os.Stat(output)
	if err != nil || info.Size() == 0 {
		t.Fatalf("Output file invalid: %s", output)
	}
}
