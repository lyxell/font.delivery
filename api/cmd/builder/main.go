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

const API_VERSION = "v2"

var ignoredFamilies = []string{
	// Ignores atma since it doesn't have a bundled license
	"atma",
	// Ignore blinker since it doesn't have a bundled license
	"blinker",
	// Ignores chathura since it doesn't have a bundled license
	"chathura",
	// Ignore dela-gothic-one since it doesn't have a bundled license
	"dela-gothic-one",
	// Ignore kulim-park since it doesn't have a bundled license
	"kulim-park",
	// Ignore mirza since it doesn't have a bundled license
	"mirza",
	// Ignore mitr since it doesn't have a bundled license
	"mitr",
	// Ignore mogra since it doesn't have a bundled license
	"mogra",
	// Ignore prata since it doesn't have a bundled license
	"prata",
	// Ignore source-serif-4 since it doesn't have a bundled license
	"source-serif-4",
	// Ignore jsMath fonts
	"jsmath-cmr10",
	"jsmath-cmex10",
	"jsmath-cmsy10",
	"jsmath-cmti10",
	"jsmath-cmmi10",
	"jsmath-cmbx10",
	// Ignore jomolhari and uchen since they use OFL-1.0 rather than OFL-1.1
	// which our code currently doesn't handle
	"uchen",
	"jomolhari",
}

func run(inputDir string, outputDir string, subsets []string) error {
	// Create needed directories
	tmpDir := "tmp"
	indexOutputDir := filepath.Join(outputDir, "api", API_VERSION)
	fontOutputDir := filepath.Join(outputDir, "api", API_VERSION, "fonts")
	licenseOutputDir := filepath.Join(outputDir, "api", API_VERSION, "licenses")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	if err := os.MkdirAll(fontOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := os.MkdirAll(licenseOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Collect metadata
	families, err := builder.CollectMetadata(inputDir, ignoredFamilies)
	if err != nil {
		return fmt.Errorf("failed to collect metadata: %w", err)
	}
	// Generate subsets JSON file
	if err := builder.GenerateSubsetsJSONFile(subsets, indexOutputDir); err != nil {
		return fmt.Errorf("failed to generate JSON file: %w", err)
	}
	// Generate index JSON file
	if err := builder.GenerateIndexJSONFile(families, subsets, indexOutputDir); err != nil {
		return fmt.Errorf("failed to generate JSON file: %w", err)
	}

	// Generate WOFF2 files
	jobs := rill.FromSlice(families, nil)
	return rill.ForEach(jobs, runtime.GOMAXPROCS(0), func(family builder.FontFamily) error {
		err := builder.GenerateLicenseFile(family, inputDir, licenseOutputDir)
		if err != nil {
			return err
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
