package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/sfhorg/font.delivery/internal/builder"
)

func run(fontPath, outputDir string, subsets []string) error {
	families, err := builder.GatherMetadata(fontPath)
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

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.GOMAXPROCS(0))

	for _, family := range families {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(f builder.FontFamily) {
			defer wg.Done()
			if err := builder.GenerateWOFF2Files(f, subsets, fontPath, outputDir); err != nil {
				log.Println("Error generating WOFF2 files:", err)
			}
			<-semaphore
		}(family)
	}
	wg.Wait()

	return nil
}

func main() {
	fontPath := flag.String("input-dir", "fonts", "Input directory containing font files")
	outputDir := flag.String("output-dir", "out", "Output directory for generated files")
	flag.Parse()

	subsets := []string{
		"latin",
		"latin-ext",
		"vietnamese",
	}

	if err := run(*fontPath, *outputDir, subsets); err != nil {
		log.Fatalf("error: %v", err)
	}
}
