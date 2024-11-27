package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/sfhorg/font.delivery/cli/internal/api"
)

func main() {
	client, err := api.NewClientWithResponses("https://font.delivery/api/v2")

	fonts, err := client.GetFontsWithResponse(context.Background())
	if err != nil {
		log.Fatalf("Error fetching fonts: %v", err)
	}

	var fontOptions []huh.Option[int]
	for i, font := range *fonts.JSON200 {
		fontOptions = append(fontOptions, huh.NewOption(font.Name, i))
	}

	var selected int
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Download a webfont").
				Description("Select a font family").
				Options(fontOptions...).
				Value(&selected).
				Height(10),
		),
	).WithTheme(huh.ThemeBase()).Run()
	if err != nil {
		log.Fatalf("Error in form: %v", err)
	}

	selectedFont := (*fonts.JSON200)[selected]

	var subsetOptions []huh.Option[string]
	for _, subset := range selectedFont.Subsets {
		subsetOptions = append(subsetOptions, huh.NewOption(string(subset), string(subset)).Selected(subset == "latin"))
	}
	var selectedSubsets []string

	var styleOptions []huh.Option[string]
	for _, style := range selectedFont.Styles {
		styleOptions = append(styleOptions, huh.NewOption(string(style), string(style)).Selected(true))
	}
	var selectedStyles []string

	var weightOptions []huh.Option[string]
	for _, weight := range selectedFont.Weights {
		weightOptions = append(weightOptions, huh.NewOption(weight, weight).Selected(true))
	}
	var selectedWeights []string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Styles").
				Options(styleOptions...).
				Value(&selectedStyles),
			huh.NewMultiSelect[string]().
				Title("Weights").
				Options(weightOptions...).
				Value(&selectedWeights),
			huh.NewMultiSelect[string]().
				Title("Subsets").
				Options(subsetOptions...).
				Value(&selectedSubsets),
		),
	).WithTheme(huh.ThemeBase()).Run()
	if err != nil {
		log.Fatalf("Error in form: %v", err)
	}

	for _, style := range selectedStyles {
		for _, subset := range selectedSubsets {
			for _, weight := range selectedWeights {
				response, err := client.DownloadFontWithResponse(
					context.Background(),
					selectedFont.Id,
					api.DownloadFontParamsSubset(subset),
					weight,
					api.DownloadFontParamsStyle(style),
				)
				if err != nil {
					log.Fatalf("Error downloading font: %v", err)
				}

				if response.StatusCode() != 200 {
					log.Fatalf("Failed to download font, HTTP status: %d", response.StatusCode())
				}

				fontFileName := fmt.Sprintf("%s_%s_%s_%s.woff2", selectedFont.Id, subset, weight, style)
				err = os.WriteFile(fontFileName, response.Body, 0o644)
				if err != nil {
					log.Fatalf("Error creating font file: %v", err)
				}

				fmt.Printf("Font downloaded and saved as %s\n", fontFileName)
			}
		}
	}
}
