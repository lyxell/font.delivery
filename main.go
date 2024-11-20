package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	sync "sync"

	"google.golang.org/protobuf/encoding/prototext"
)

type Font struct {
	Name       string `json:"name"`
	Style      string `json:"style"`
	Weight     int    `json:"weight"`
	Filename   string `json:"filename"`
	PostScript string `json:"post_script_name"`
	FullName   string `json:"full_name"`
	Copyright  string `json:"copyright"`
}

type Axis struct {
	Tag      string  `json:"tag"`
	MinValue float32 `json:"min_value"`
	MaxValue float32 `json:"max_value"`
}

type FontFamily struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Designer string   `json:"designer"`
	License  string   `json:"license"`
	Category []string `json:"category"`
	Fonts    []Font   `json:"fonts"`
	Subsets  []string `json:"subsets"`
	Axes     []Axis   `json:"axes"`
	Minisite string   `json:"minisite_url"`
}

func ReadProtoFromText(path string) (*FamilyProto, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	protoInstance := &FamilyProto{}
	if err := prototext.Unmarshal(data, protoInstance); err != nil {
		return nil, err
	}
	return protoInstance, nil
}

// GatherAllMetadata walks through the directory and gathers metadata from
// Protobuf files.
func GatherAllMetadata(rootDir string) ([]FontFamily, error) {
	var allMetadata []FontFamily

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "METADATA.pb" {
			familyData, err := ReadProtoFromText(path)
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
				family.Fonts = append(family.Fonts, Font{
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
				family.Axes = append(family.Axes, Axis{
					Tag:      axisProto.GetTag(),
					MinValue: axisProto.GetMinValue(),
					MaxValue: axisProto.GetMaxValue(),
				})
			}

			allMetadata = append(allMetadata, family)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allMetadata, nil
}

func readUnicodeRanges(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var ranges []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ranges = append(ranges, "U+"+line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(ranges, ", "), nil
}

func getFontWeight(fontData FontFamily, font Font) string {
	for _, axis := range fontData.Axes {
		if axis.Tag == "wght" {
			return fmt.Sprintf("%v %v", axis.MinValue, axis.MaxValue)
		}
	}
	return fmt.Sprintf("%d", font.Weight)
}

func generateCSS(fontData FontFamily, subsets []string) string {
	var cssOutput string
	for _, subset := range subsets {
		unicodeRanges, err := readUnicodeRanges("ranges/" + subset + ".txt")
		if err != nil {
			log.Fatalf("Failed to read Unicode ranges: %v", err)
		}

		for _, font := range fontData.Fonts {
			fontWeight := getFontWeight(fontData, font)
			cssOutput += fmt.Sprintf(`@font-face {
	font-family: "%s";
	font-style: %s;
	font-weight: %s;
	font-display: swap;
	src: url('%s_%s_%s_%s.woff2') format('woff2');
	unicode-range: %s;
}
`, fontData.Name, font.Style, fontWeight, fontData.Id, subset, strings.Replace(fontWeight, " ", "-", 1), font.Style, unicodeRanges)
		}
	}

	cssOutput += fmt.Sprintf(`.font-%s {
  font-family: "%s";
}
`, fontData.Id, fontData.Name)
	return cssOutput
}

func getLicenseDirName(license string) string {
	switch strings.ToLower(license) {
	case "apache2":
		return "apache"
	default:
		return strings.ToLower(license)
	}
}

func generateCSSFiles(families []FontFamily) {
	for _, f := range families {
		css := generateCSS(f, []string{"latin"})
		err := os.WriteFile("out/"+f.Id+".css", []byte(css), 0o644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func generateWoff2Files(family FontFamily, fontPath string, subsets []string) error {
	licenseDir := getLicenseDirName(family.License)
	for _, font := range family.Fonts {
		inputPath := fmt.Sprintf("%s/%s/%s/%s", fontPath, licenseDir, strings.ToLower(strings.ReplaceAll(family.Name, " ", "")), font.Filename)

		for _, subset := range subsets {
			if !slices.Contains(family.Subsets, subset) {
				continue
			}
			unicodeFilePath := "ranges/" + subset + ".txt"
			tempSubsetPath := fmt.Sprintf("temp/%s_%s.subset.ttf", family.Id, subset)
			outputPath := fmt.Sprintf("out/%s_%s_%s_%s.woff2", family.Id, subset, strings.Replace(getFontWeight(family, font), " ", "-", -1), font.Style)

			cmdSubset := exec.Command("hb-subset", "--unicodes-file="+unicodeFilePath, "--output-file="+tempSubsetPath, inputPath)
			if err := cmdSubset.Run(); err != nil {
				return fmt.Errorf("error subsetting font %s for subset %s: %w", font.Name, subset, err)
			}
			cmdCompress := exec.Command("woff2_compress", tempSubsetPath)
			if err := cmdCompress.Run(); err != nil {
				return fmt.Errorf("error compressing to WOFF2 for font %s, subset %s: %w", font.Name, subset, err)
			}

			tempWoff2Path := tempSubsetPath[:len(tempSubsetPath)-len(".ttf")] + ".woff2"
			if err := os.Rename(tempWoff2Path, outputPath); err != nil {
				return fmt.Errorf("error moving WOFF2 file to output directory for font %s, subset %s: %w", font.Name, subset, err)
			}
		}
	}
	return nil
}

func main() {
	metadataRoot := "/home/user/projects/fonts/ofl"
	fontPath := "/home/user/projects/fonts"
	subsets := []string{"latin"}

	families, err := GatherAllMetadata(metadataRoot)
	if err != nil {
		log.Fatalf("Failed to gather metadata: %v", err)
	}

	os.MkdirAll("temp", os.ModePerm)
	os.MkdirAll("out", os.ModePerm)

	generateCSSFiles(families)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.GOMAXPROCS(0))

	for _, family := range families {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(f FontFamily) {
			defer wg.Done()
			if err := generateWoff2Files(f, fontPath, subsets); err != nil {
				log.Println("Error generating WOFF2 files:", err)
			}
			<-semaphore
		}(family)
	}
	wg.Wait()
}
