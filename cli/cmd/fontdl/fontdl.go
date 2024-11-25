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
	client, err := api.NewClientWithResponses("https://font.delivery/api/v1")

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

	fontDetails, err := client.GetFontByIDWithResponse(context.Background(), selectedFont.Id)
	if err != nil {
		log.Fatalf("Error fetching font details: %v", err)
	}

	var subsetOptions []huh.Option[string]
	for _, subset := range fontDetails.JSON200.Subsets {
		subsetOptions = append(subsetOptions, huh.NewOption(string(subset), string(subset)))
	}
	var selectedSubset string

	var styleOptions []huh.Option[string]
	for _, style := range fontDetails.JSON200.Styles {
		styleOptions = append(styleOptions, huh.NewOption(string(style), string(style)))
	}
	var selectedStyle string

	var weightOptions []huh.Option[string]
	for _, weight := range fontDetails.JSON200.Weights {
		weightOptions = append(weightOptions, huh.NewOption(weight, weight))
	}
	var selectedWeight string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a subset").
				Options(subsetOptions...).
				Value(&selectedSubset),
			huh.NewSelect[string]().
				Title("Select a style").
				Options(styleOptions...).
				Value(&selectedStyle),
			huh.NewSelect[string]().
				Title("Select a weight").
				Options(weightOptions...).
				Value(&selectedWeight),
		),
	).WithTheme(huh.ThemeBase()).Run()
	if err != nil {
		log.Fatalf("Error in form: %v", err)
	}

	response, err := client.DownloadFontWithResponse(
		context.Background(),
		selectedFont.Id,
		api.DownloadFontParamsSubset(selectedSubset),
		selectedWeight,
		api.DownloadFontParamsStyle(selectedStyle),
	)
	if err != nil {
		log.Fatalf("Error downloading font: %v", err)
	}

	if response.StatusCode() != 200 {
		log.Fatalf("Failed to download font, HTTP status: %d", response.StatusCode())
	}

	fontFileName := fmt.Sprintf("%s_%s_%s_%s.woff2", selectedFont.Id, selectedSubset, selectedWeight, selectedStyle)
	err = os.WriteFile(fontFileName, response.Body, 0o644)
	if err != nil {
		log.Fatalf("Error creating font file: %v", err)
	}

	fmt.Printf("Font downloaded and saved as %s\n", fontFileName)
}
