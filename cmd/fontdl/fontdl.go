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

	fmt.Println(selectedFont.Name)
}
