package pdf_utils

import (
	"github.com/hamzaelmarjani/pdf-utils-go/utils"
)

/// PDFUtils is a high-level utility for PDF operations.
type PDFUtils struct{}

/// NewPDFUtils creates a new instance of PDFUtils.
func NewPDFUtils() *PDFUtils {
	return &PDFUtils{}
}

/// Combine multiple PDF files into one.
func (u *PDFUtils) Combine(paths []string, output string) (string, error) {
	return utils.Combine(paths, output)
}

/// Split a PDF file into individual pages.
func (u *PDFUtils) Split(path string) ([]string, error) {
	return utils.Split(path)
}

/// Compress a PDF file.
func (u *PDFUtils) Compress(path string, percentage float32) (string, error) {
	return utils.Compress(path, percentage)
}

/// Convert a file to PDF.
func (u *PDFUtils) ConvertToPDF(path string) (string, error) {
	return utils.ToPDF(path)
}

/// Convert a PDF to another format.
func (u *PDFUtils) ReverseFromPDF(path string, format string) (string, error) {
	return utils.FromPDF(path, format)
}
