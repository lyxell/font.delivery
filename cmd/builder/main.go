package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/destel/rill"
	"github.com/sfhorg/font.delivery/internal/builder"
)

func run(inputDir string, outputDir string, subsets []string) error {
	families, err := builder.GatherMetadata(inputDir)
	if err != nil {
		return fmt.Errorf("failed to gather metadata: %w", err)
	}
	if err := os.MkdirAll("temp", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := builder.GenerateJSONFiles(families, subsets, outputDir); err != nil {
		return fmt.Errorf("failed to generate JSON files: %w", err)
	}

	if err := builder.GenerateCSSFiles(families, subsets, outputDir); err != nil {
		return fmt.Errorf("failed to generate CSS files: %w", err)
	}

	jobs := rill.FromSlice(families, nil)

	// Create needed directories
	fontOutputDir := filepath.Join(outputDir, "api", "v1", "download")
	err = os.MkdirAll(fontOutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	return rill.ForEach(jobs, runtime.GOMAXPROCS(0), func(f builder.FontFamily) error {
		return builder.GenerateWOFF2Files(f, subsets, inputDir, fontOutputDir)
	})
}

func main() {
	inputDir := flag.String("input-dir", "fonts", "Input directory containing font files")
	outputDir := flag.String("output-dir", "out", "Output directory for generated files")
	flag.Parse()

	subsets := []string{
		"latin",
		"latin-ext",
		"vietnamese",
	}

	if err := run(*inputDir, *outputDir, subsets); err != nil {
		log.Fatalf("error: %v", err)
	}
}
