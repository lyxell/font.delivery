package main

import (
	"context"
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/sfhorg/fontdelivery/internal/api"
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

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a subset").
				Options(subsetOptions...).
				Value(&selectedSubset).
				Height(5),
			huh.NewSelect[string]().
				Title("Select a style").
				Options(styleOptions...).
				Value(&selectedStyle).
				Height(5),
		),
	).WithTheme(huh.ThemeBase()).Run()
	if err != nil {
		log.Fatalf("Error in form: %v", err)
	}

	fmt.Printf("Selected Font: %s\n", selectedFont.Name)
	fmt.Printf("Selected Subset: %s\n", selectedSubset)
	fmt.Printf("Selected Style: %s\n", selectedStyle)
}
