package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/lyxell/font.delivery/cli/internal/api"
)

// generateFontFaceCSS generates the @font-face CSS rule for a font
func generateFontFaceCSS(
	fontName, fontID, subset, weight, style, unicodeRange string,
) string {
	url := fmt.Sprintf("%s_%s_%s_%s.woff2", fontID, subset, weight, style)
	return strings.TrimSpace(fmt.Sprintf(`
@font-face {
  font-family: '%s';
  font-style: %s;
  font-weight: %s;
  src: url('%s') format('woff2');
  unicode-range: %s;
}
`, fontName, style, strings.Replace(weight, "-", " ", 1), url, unicodeRange))
}

func run() error {
	client, err := api.NewClientWithResponses("https://font.delivery/api/v2")
	if err != nil {
		return fmt.Errorf("creating API client: %w", err)
	}

	fonts, err := client.GetFontsWithResponse(context.Background())
	if err != nil {
		return fmt.Errorf("fetching fonts: %w", err)
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
		return fmt.Errorf("form selection error: %w", err)
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
		return fmt.Errorf("form multi-select error: %w", err)
	}

	// Fetch subsets
	subsetsResponse, err := client.GetSubsetsWithResponse(context.Background())
	if err != nil || subsetsResponse.JSON200 == nil {
		return fmt.Errorf("fetching subsets: %w", err)
	}

	// Build a map of subset to unicode ranges
	subsetRanges := make(map[string]string)
	for _, subset := range *subsetsResponse.JSON200 {
		subsetRanges[string(subset.Subset)] = subset.Ranges
	}

	// Download license file
	licenseResponse, err := client.DownloadLicenseWithResponse(context.Background(), selectedFont.Id)
	if err != nil {
		return fmt.Errorf("downloading license: %w", err)
	}

	if licenseResponse.StatusCode() != 200 {
		return fmt.Errorf("failed to download license, HTTP status: %d", licenseResponse.StatusCode())
	}

	licenseFileName := fmt.Sprintf("%s-LICENSE.txt", selectedFont.Id)
	err = os.WriteFile(licenseFileName, licenseResponse.Body, 0o644)
	if err != nil {
		return fmt.Errorf("writing license file: %w", err)
	}
	fmt.Printf("License file downloaded and saved as %s\n", licenseFileName)

	var cssContent strings.Builder
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
					return fmt.Errorf("downloading font: %w", err)
				}

				if response.StatusCode() != 200 {
					return fmt.Errorf("failed to download font, HTTP status: %d", response.StatusCode())
				}

				fontFileName := fmt.Sprintf("%s_%s_%s_%s.woff2", selectedFont.Id, subset, weight, style)
				err = os.WriteFile(fontFileName, response.Body, 0o644)
				if err != nil {
					return fmt.Errorf("writing font file: %w", err)
				}

				fmt.Printf("Font downloaded and saved as %s\n", fontFileName)
				cssContent.WriteString(generateFontFaceCSS(
					selectedFont.Name,
					selectedFont.Id,
					string(subset),
					weight,
					string(style),
					subsetRanges[string(subset)],
				))
				cssContent.WriteString("\n")
			}
		}
	}

	cssFileName := selectedFont.Id + ".css"
	err = os.WriteFile(cssFileName, []byte(cssContent.String()), 0o644)
	if err != nil {
		return fmt.Errorf("writing CSS file: %w", err)
	}
	fmt.Printf("CSS file generated: %s\n", cssFileName)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
