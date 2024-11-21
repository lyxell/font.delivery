package main

import (
	"encoding/json"
	"testing"
)

func TestGenerateCSSInter(t *testing.T) {
	testJSON := `{
    "id": "inter",
    "name": "Inter",
    "designer": "Rasmus Andersson",
    "license": "OFL",
    "category": [
      "SANS_SERIF"
    ],
    "fonts": [
      {
        "name": "Inter",
        "style": "normal",
        "weight": 400,
        "filename": "Inter[slnt,wght].ttf",
        "post_script_name": "Inter-Regular",
        "full_name": "Inter Regular",
        "copyright": "Copyright 2020 The Inter Project Authors (https://github.com/rsms/inter)"
      }
    ],
    "subsets": [
      "cyrillic",
      "cyrillic-ext",
      "greek",
      "greek-ext",
      "latin",
      "latin-ext",
      "menu",
      "vietnamese"
    ],
    "axes": [
      {
        "tag": "slnt",
        "min_value": -10.0,
        "max_value": 0.0
      },
      {
        "tag": "wght",
        "min_value": 100.0,
        "max_value": 900.0
      }
    ],
    "source": {
      "repository_url": "https://github.com/rsms/inter"
    },
    "minisite_url": "https://rsms.me/inter/"
  }`

	var testData FontFamily

	if err := json.Unmarshal([]byte(testJSON), &testData); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	generatedCSS, err := generateCSS(testData, []string{"latin"})
	if err != nil {
		t.Fatalf("Expected err to be nil: %v", err)
	}

	expectedCSS := `@font-face {
	font-family: "Inter";
	font-style: normal;
	font-weight: 100 900;
	font-display: swap;
	src: url('inter_latin_100-900_normal.woff2') format('woff2');
	unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}
.font-inter {
  font-family: "Inter";
}
`

	if generatedCSS != expectedCSS {
		t.Errorf("Generated CSS does not match expected CSS.\nExpected:\n%s\nGot:\n%s", expectedCSS, generatedCSS)
	}
}

func TestGenerateCSSJetBrainsMono(t *testing.T) {
	testJSON := `{
    "id": "jetbrains-mono",
    "name": "JetBrains Mono",
    "designer": "JetBrains, Philipp Nurullin, Konstantin Bulenkov",
    "license": "OFL",
    "category": [
      "MONOSPACE"
    ],
    "fonts": [
      {
        "name": "JetBrains Mono",
        "style": "normal",
        "weight": 400,
        "filename": "JetBrainsMono[wght].ttf",
        "post_script_name": "JetBrainsMono-Regular",
        "full_name": "JetBrains Mono Regular",
        "copyright": "Copyright 2020 The JetBrains Mono Project Authors (https://github.com/JetBrains/JetBrainsMono)"
      },
      {
        "name": "JetBrains Mono",
        "style": "italic",
        "weight": 400,
        "filename": "JetBrainsMono-Italic[wght].ttf",
        "post_script_name": "JetBrainsMono-Italic",
        "full_name": "JetBrains Mono Italic",
        "copyright": "Copyright 2020 The JetBrains Mono Project Authors (https://github.com/JetBrains/JetBrainsMono)"
      }
    ],
    "subsets": [
      "cyrillic",
      "cyrillic-ext",
      "greek",
      "latin",
      "latin-ext",
      "menu",
      "vietnamese"
    ],
    "axes": [
      {
        "tag": "wght",
        "min_value": 100.0,
        "max_value": 800.0
      }
    ],
    "minisite_url": "https://www.jetbrains.com/lp/mono/"
  }`

	var testData FontFamily

	if err := json.Unmarshal([]byte(testJSON), &testData); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	generatedCSS, err := generateCSS(testData, []string{"latin"})
	if err != nil {
		t.Fatalf("Expected err to be nil: %v", err)
	}

	expectedCSS := `@font-face {
	font-family: "JetBrains Mono";
	font-style: normal;
	font-weight: 100 800;
	font-display: swap;
	src: url('jetbrains-mono_latin_100-800_normal.woff2') format('woff2');
	unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}
@font-face {
	font-family: "JetBrains Mono";
	font-style: italic;
	font-weight: 100 800;
	font-display: swap;
	src: url('jetbrains-mono_latin_100-800_italic.woff2') format('woff2');
	unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}
.font-jetbrains-mono {
  font-family: "JetBrains Mono";
}
`

	if generatedCSS != expectedCSS {
		t.Errorf("Generated CSS does not match expected CSS.\nExpected:\n%s\nGot:\n%s", expectedCSS, generatedCSS)
	}
}

func TestGenerateCSSJoan(t *testing.T) {
	testJSON := `{
    "id": "joan",
    "name": "Joan",
    "designer": "Paolo Biagini",
    "license": "OFL",
    "category": [
      "SERIF"
    ],
    "fonts": [
      {
        "name": "Joan",
        "style": "normal",
        "weight": 400,
        "filename": "Joan-Regular.ttf",
        "post_script_name": "Joan-Regular",
        "full_name": "Joan Regular",
        "copyright": "Copyright 2021 The Joan Project Authors (https://github.com/PaoloBiagini/Joan)"
      }
    ],
    "subsets": [
      "latin",
      "latin-ext",
      "menu"
    ],
    "axes": [],
    "source": {
      "repository_url": "https://github.com/PaoloBiagini/Joan",
      "commit": "981cb73299f7d9164eedcb647e57fb34c9dc1139",
      "files": [
        {
          "source_file": "OFL.txt",
          "dest_file": "OFL.txt"
        },
        {
          "source_file": "fonts/ttf/Joan-Regular.ttf",
          "dest_file": "Joan-Regular.ttf"
        }
      ],
      "branch": "main"
    },
    "minisite_url": null
  }`

	var testData FontFamily

	if err := json.Unmarshal([]byte(testJSON), &testData); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	expectedCSS := `@font-face {
	font-family: "Joan";
	font-style: normal;
	font-weight: 400;
	font-display: swap;
	src: url('joan_latin_400_normal.woff2') format('woff2');
	unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}
.font-joan {
  font-family: "Joan";
}
`

	generatedCSS, err := generateCSS(testData, []string{"latin"})
	if err != nil {
		t.Fatalf("Expected err to be nil: %v", err)
	}

	if generatedCSS != expectedCSS {
		t.Errorf("Generated CSS does not match expected CSS.\nExpected:\n%s\nGot:\n%s", expectedCSS, generatedCSS)
	}
}
