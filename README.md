# PDF Utils Library (Go)

A comprehensive Go library for PDF manipulation, conversion, and more.

## Features

- **Combiner**: Merge multiple PDF files into one.
- **Splitter**: Split a PDF into individual pages.
- **Compressor**: Optimize PDF files to reduce size.
- **Converter**: Convert various formats (Text, Images, HTML) to PDF.
- **Reverser**: Convert PDF to Text, HTML, or Images.

## Requirements

- **HTML Conversion & PDF to Image**: Requires a Chrome or Chromium installation on the system as it uses `chromedp`.
- **Platform Support**: Works on Windows, macOS, and Linux.

## Installation

```bash
go get github.com/hamzaelmarjani/pdf-utils-go
```

## Usage

### Initialize the Utility

```go
import "github.com/hamzaelmarjani/pdf-utils-go"

utils := pdf_utils.NewPDFUtils()
```

### 1. Combine PDFs

```go
paths := []string{"file1.pdf", "file2.pdf"}
output, err := utils.Combine(paths, "")
if err == nil {
    fmt.Printf("Combined PDF saved at: %s\n", output)
}
```

### 2. Split PDF

```go
pages, err := utils.Split("large_file.pdf")
if err == nil {
    for i, pagePath := range pages {
        fmt.Printf("Page %d saved at: %s\n", i+1, pagePath)
    }
}
```

### 3. Compress PDF

```go
compressedPath, err := utils.Compress("original.pdf", 0.5)
if err == nil {
    fmt.Printf("Compressed PDF saved at: %s\n", compressedPath)
}
```

### 4. Convert File to PDF

Supported formats: `txt`, `png`, `jpg`, `jpeg`, `html`.

```go
pdfPath, err := utils.ConvertToPDF("image.jpeg")
if err == nil {
    fmt.Printf("Generated PDF: %s\n", pdfPath)
}
```

### 5. Convert PDF to Other Formats

Supported formats: `text`, `html`, `png`, `jpg`, `jpeg`.

```go
textPath, err := utils.ReverseFromPDF("document.pdf", "text")
if err == nil {
    fmt.Printf("Extracted text saved at: %s\n", textPath)
}
```

## Test

1. Clone the repository and enter the project folder:

```bash
git clone https://github.com/hamzaelmarjani/pdf-utils-go.git pdf-utils-go && cd pdf-utils-go
```

2. Create a new `tests_data` folder in the project root, and set the source files there.
3. Run the tests:

```bash
go test
```

4. Check the results inside the `./dist` folder.

## License

MIT
