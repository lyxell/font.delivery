package builder

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sfhorg/fontdelivery/internal/subsetting"
	"google.golang.org/protobuf/encoding/prototext"
)

type FontFamilyFont struct {
	Name       string `json:"name"`
	Style      string `json:"style"`
	Weight     int    `json:"weight"`
	Filename   string `json:"filename"`
	PostScript string `json:"post_script_name"`
	FullName   string `json:"full_name"`
	Copyright  string `json:"copyright"`
}

type FontFamilyAxis struct {
	Tag      string  `json:"tag"`
	MinValue float32 `json:"min_value"`
	MaxValue float32 `json:"max_value"`
}

type FontFamily struct {
	Id       string           `json:"id"`
	Name     string           `json:"name"`
	Designer string           `json:"designer"`
	License  string           `json:"license"`
	Category []string         `json:"category"`
	Fonts    []FontFamilyFont `json:"fonts"`
	Subsets  []string         `json:"subsets"`
	Axes     []FontFamilyAxis `json:"axes"`
	Minisite string           `json:"minisite_url"`
}

// parseMetadataProtobuf parses a METADATA.pb file into a FamilyProto struct.
func parseMetadataProtobuf(path string) (*FamilyProto, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var protoInstance FamilyProto
	if err := prototext.Unmarshal(data, &protoInstance); err != nil {
		return nil, err
	}
	return &protoInstance, nil
}

// GatherMetadata walks through the supplied directory and gathers metadata
// from all METADATA.pb files it stumbles upon.
//
// The slice of metadata will be sorted by the name of the font family.
func GatherMetadata(rootDir string) ([]FontFamily, error) {
	var metadata []FontFamily
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "METADATA.pb" {
			familyData, err := parseMetadataProtobuf(path)
			if err != nil {
				return err
			}
			family := FontFamily{
				Id:       strings.ToLower(strings.ReplaceAll(familyData.GetName(), " ", "-")),
				Name:     familyData.GetName(),
				Designer: familyData.GetDesigner(),
				License:  familyData.GetLicense(),
				Category: familyData.GetCategory(),
				Subsets:  familyData.GetSubsets(),
				Minisite: familyData.GetMinisiteUrl(),
			}
			for _, fontProto := range familyData.GetFonts() {
				family.Fonts = append(family.Fonts, FontFamilyFont{
					Name:       fontProto.GetName(),
					Style:      fontProto.GetStyle(),
					Weight:     int(fontProto.GetWeight()),
					Filename:   fontProto.GetFilename(),
					PostScript: fontProto.GetPostScriptName(),
					FullName:   fontProto.GetFullName(),
					Copyright:  fontProto.GetCopyright(),
				})
			}
			for _, axisProto := range familyData.GetAxes() {
				family.Axes = append(family.Axes, FontFamilyAxis{
					Tag:      axisProto.GetTag(),
					MinValue: axisProto.GetMinValue(),
					MaxValue: axisProto.GetMaxValue(),
				})
			}
			metadata = append(metadata, family)
		}
		return nil
	})
	slices.SortFunc(metadata, func(a, b FontFamily) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return metadata, err
}

func getFontWeight(fontData FontFamily, font FontFamilyFont) string {
	for _, axis := range fontData.Axes {
		if axis.Tag == "wght" {
			return fmt.Sprintf("%v %v", axis.MinValue, axis.MaxValue)
		}
	}
	return fmt.Sprintf("%d", font.Weight)
}

func generateCSS(fontData FontFamily, subsets []string) string {
	var cssOutput strings.Builder

	fontFaceTemplate := `@font-face {
	font-family: "{family}";
	font-style: {style};
	font-weight: {weight};
	font-display: swap;
	src: url('{fileName}.woff2') format('woff2');
	unicode-range: {ranges};
}
`
	for _, subset := range subsets {
		ranges := subsetting.BuildCSSString(subset)
		for _, font := range fontData.Fonts {
			fontWeight := getFontWeight(fontData, font)
			fileName := strings.Join([]string{
				fontData.Id,
				subset,
				strings.Replace(fontWeight, " ", "-", 1),
				font.Style,
			}, "_")
			replacer := strings.NewReplacer(
				"{family}", fontData.Name,
				"{style}", font.Style,
				"{weight}", fontWeight,
				"{fileName}", fileName,
				"{ranges}", ranges,
			)
			cssOutput.WriteString(replacer.Replace(fontFaceTemplate))
		}
	}

	cssOutput.WriteString(fmt.Sprintf(`.font-%s {
  font-family: "%s";
}
`, fontData.Id, fontData.Name))

	return cssOutput.String()
}

func getLicenseDirName(license string) string {
	switch strings.ToLower(license) {
	case "apache2":
		return "apache"
	default:
		return strings.ToLower(license)
	}
}

func GenerateCSSFiles(families []FontFamily, subsets []string, outputDir string) error {
	for _, f := range families {
		css := generateCSS(f, subsets)
		err := os.WriteFile(filepath.Join(outputDir, f.Id+".css"), []byte(css), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeHarfbuzzRangeFile(path string, subset string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(subsetting.BuildHarfbuzzString(subset)))
	if err != nil {
		return err
	}
	return nil
}

func GenerateWOFF2Files(family FontFamily, subsets []string, fontPath string, outputDir string) error {
	const tempPath = "temp"

	for _, subset := range subsets {
		// We add the family.Id here to avoid race conditions where goroutines
		// could overwrite the files of other goroutines
		unicodeRangesPath := filepath.Join(tempPath, fmt.Sprintf("range-%s-%s.txt", family.Id, subset))
		writeHarfbuzzRangeFile(unicodeRangesPath, subset)
	}

	licenseDir := getLicenseDirName(family.License)
	for _, font := range family.Fonts {
		// inputPath is where we find the original .tff-file
		inputPath := filepath.Join(
			fontPath,
			licenseDir,
			strings.ToLower(strings.ReplaceAll(family.Name, " ", "")),
			font.Filename,
		)
		fmt.Println("Processing", inputPath)
		for _, subset := range subsets {
			if !slices.Contains(family.Subsets, subset) {
				continue
			}

			// unicodeRangesPath is where harfbuzz reads the unicode ranges for subsetting from
			unicodeRangesPath := filepath.Join(tempPath, fmt.Sprintf("range-%s-%s.txt", family.Id, subset))

			// tempSubsetPath is where the intermediary subsetted .ttf-file will be written to
			tempSubsetPath := filepath.Join(tempPath, fmt.Sprintf("%s_%s.subset.ttf", family.Id, subset))

			// outputPath is where the final .woff2-file will be written to
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s_%s_%s.woff2", family.Id, subset, strings.Replace(getFontWeight(family, font), " ", "-", -1), font.Style))

			// Perform subsetting
			cmdSubset := exec.Command("hb-subset", "--unicodes-file="+unicodeRangesPath, "--output-file="+tempSubsetPath, inputPath)
			if err := cmdSubset.Run(); err != nil {
				return fmt.Errorf("error subsetting font %s for subset %s: %w", font.Name, subset, err)
			}

			// Generate woff2-file
			cmdCompress := exec.Command("woff2_compress", tempSubsetPath)
			if err := cmdCompress.Run(); err != nil {
				return fmt.Errorf("error compressing to WOFF2 for font %s, subset %s: %w", font.Name, subset, err)
			}

			// Move file to final destination
			tempWoff2Path := strings.TrimSuffix(tempSubsetPath, ".ttf") + ".woff2"
			if err := os.Rename(tempWoff2Path, outputPath); err != nil {
				return fmt.Errorf("error moving WOFF2 file to output directory for font %s, subset %s: %w", font.Name, subset, err)
			}
		}
	}
	return nil
}

func GenerateJSONFiles(families []FontFamily, subsets []string, outputDir string) error {
	var apiData []map[string]string
	for _, font := range families {
		subsetsIntersect := false
		for _, s := range subsets {
			if slices.Contains(font.Subsets, s) {
				subsetsIntersect = true
				break
			}
		}
		// Skip fonts that do not have any renderable subsets
		if !subsetsIntersect {
			continue
		}
		apiData = append(apiData, map[string]string{
			"id":       font.Id,
			"name":     font.Name,
			"designer": font.Designer,
			"css":      fmt.Sprintf("/%s.css", font.Id),
		})
	}

	apiDir := filepath.Join(outputDir, "api", "v1")
	err := os.MkdirAll(apiDir, os.ModePerm)
	if err != nil {
		return err
	}

	apiDataBytes, err := json.MarshalIndent(apiData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(apiDir, "fonts.json"), apiDataBytes, 0o644)
}
