package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFontWeightsWithSampleData(t *testing.T) {
	family := FontFamily{
		Id:       "alegreya-sans",
		Name:     "Alegreya Sans",
		Designer: "Juan Pablo del Peral, Huerta Tipográfica",
		License:  "OFL",
		Category: []string{"SANS_SERIF"},
		Fonts: []FontFamilyFont{
			{Weight: 100, Style: "normal"},
			{Weight: 100, Style: "italic"},
			{Weight: 300, Style: "normal"},
			{Weight: 300, Style: "italic"},
			{Weight: 400, Style: "normal"},
			{Weight: 400, Style: "italic"},
			{Weight: 500, Style: "normal"},
			{Weight: 500, Style: "italic"},
			{Weight: 700, Style: "normal"},
			{Weight: 700, Style: "italic"},
			{Weight: 800, Style: "normal"},
			{Weight: 800, Style: "italic"},
			{Weight: 900, Style: "normal"},
			{Weight: 900, Style: "italic"},
		},
		Subsets:  []string{"cyrillic", "cyrillic-ext", "greek", "greek-ext", "latin", "latin-ext", "menu", "vietnamese"},
		Axes:     nil,
		Minisite: "https://huertatipografica.com/en/fonts/alegreya-sans-ht",
	}

	expectedWeights := []string{
		"100", "300", "400", "500", "700", "800", "900",
	}

	actualWeights := getFontWeights(family)

	assert.Equal(t, actualWeights, expectedWeights)
}

func TestGetFontStylesWithSampleData(t *testing.T) {
	family := FontFamily{
		Id:       "alegreya-sans",
		Name:     "Alegreya Sans",
		Designer: "Juan Pablo del Peral, Huerta Tipográfica",
		License:  "OFL",
		Category: []string{"SANS_SERIF"},
		Fonts: []FontFamilyFont{
			{Weight: 100, Style: "normal"},
			{Weight: 100, Style: "italic"},
			{Weight: 300, Style: "normal"},
			{Weight: 300, Style: "italic"},
			{Weight: 400, Style: "normal"},
			{Weight: 400, Style: "italic"},
			{Weight: 500, Style: "normal"},
			{Weight: 500, Style: "italic"},
			{Weight: 700, Style: "normal"},
			{Weight: 700, Style: "italic"},
			{Weight: 800, Style: "normal"},
			{Weight: 800, Style: "italic"},
			{Weight: 900, Style: "normal"},
			{Weight: 900, Style: "italic"},
		},
		Subsets:  []string{"cyrillic", "cyrillic-ext", "greek", "greek-ext", "latin", "latin-ext", "menu", "vietnamese"},
		Axes:     nil,
		Minisite: "https://huertatipografica.com/en/fonts/alegreya-sans-ht",
	}

	expectedStyles := []string{
		"normal", "italic",
	}

	actualStyles := getFontStyles(family)

	assert.Equal(t, actualStyles, expectedStyles)
}

