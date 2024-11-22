package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/huh"
)

type Font struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FontDetails struct {
	ID      string   `json:"id"`
	Subsets []string `json:"subsets"`
	Styles  []string `json:"styles"`
}

func fetchFonts(url string) ([]Font, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch fonts: %s", resp.Status)
	}

	var fonts []Font
	err = json.NewDecoder(resp.Body).Decode(&fonts)
	if err != nil {
		return nil, err
	}

	return fonts, nil
}

func fetchFontDetails(url string) (*FontDetails, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch font details: %s", resp.Status)
	}

	var details FontDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func main() {
	var theme *huh.Theme = huh.ThemeBase()

	fonts, err := fetchFonts("https://font.delivery/api/v1/fonts.json")
	if err != nil {
		log.Fatal(err)
	}

	var fontOptions []huh.Option[int]

	for i, font := range fonts {
		fontOptions = append(fontOptions, huh.NewOption(font.Name, i))
	}
	var selected int

	fmt.Println()
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Download a webfont").
				Description(" ").
				Options(fontOptions...).
				Value(&selected).
				Height(10),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}

	selectedFont := fonts[selected]

	// Fetch font details
	fontDetails, err := fetchFontDetails(fmt.Sprintf("https://font.delivery/api/v1/fonts/%s.json", selectedFont.ID))
	if err != nil {
		log.Fatal(err)
	}

	var subsetOptions []huh.Option[string]
	for _, subset := range fontDetails.Subsets {
		subsetOptions = append(subsetOptions, huh.NewOption(subset, subset))
	}
	var selectedSubset string

	var styleOptions []huh.Option[string]
	for _, style := range fontDetails.Styles {
		styleOptions = append(styleOptions, huh.NewOption(style, style))
	}
	var selectedStyle string

	fmt.Println()
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
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the selected font, subset, and style
	fmt.Printf("Selected Font: %s\n", selectedFont.Name)
	fmt.Printf("Selected Subset: %s\n", selectedSubset)
	fmt.Printf("Selected Style: %s\n", selectedStyle)
}
