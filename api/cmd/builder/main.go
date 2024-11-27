package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/destel/rill"
	"github.com/sfhorg/font.delivery/api/internal/builder"
)

func run(inputDir string, outputDir string, subsets []string) error {
	// Create needed directories
	tmpDir := "tmp"
	indexOutputDir := filepath.Join(outputDir, "api", "v1")
	fontOutputDir := filepath.Join(outputDir, "api", "v1", "download")
	jsonOutputDir := filepath.Join(outputDir, "api", "v1", "fonts")
	cssOutputDir := filepath.Join(outputDir, "api", "v1", "css")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	if err := os.MkdirAll(fontOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := os.MkdirAll(jsonOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create json output directory: %w", err)
	}
	if err := os.MkdirAll(cssOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create css output directory: %w", err)
	}

	// Collect metadata
	families, err := builder.CollectMetadata(inputDir)
	if err != nil {
		return fmt.Errorf("failed to collect metadata: %w", err)
	}
	// Generate index JSON file
	if err := builder.GenerateIndexJSONFile(families, subsets, indexOutputDir); err != nil {
		return fmt.Errorf("failed to generate JSON file: %w", err)
	}

	// Generate files
	jobs := rill.FromSlice(families, nil)
	return rill.ForEach(jobs, runtime.GOMAXPROCS(0), func(family builder.FontFamily) error {
		fmt.Println("Building", family.Name)
		if err := builder.GenerateFamilyJSONFile(family, subsets, jsonOutputDir); err != nil {
			return fmt.Errorf("failed to generate JSON file: %w", err)
		}
		if err := builder.GenerateFamilyCSSFiles(family, subsets, cssOutputDir); err != nil {
			return fmt.Errorf("failed to generate CSS files: %w", err)
		}
		return builder.GenerateWOFF2Files(family, subsets, inputDir, fontOutputDir, tmpDir)
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
		"cyrillic",
		"cyrillic-ext",
		"hebrew",
		"greek",
		"greek-ext",
	}

	if err := run(*inputDir, *outputDir, subsets); err != nil {
		log.Fatalf("error: %v", err)
	}
}
