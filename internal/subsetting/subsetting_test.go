package subsetting_test

import (
	"testing"

	"github.com/sfhorg/fontdelivery/internal/subsetting"
	"github.com/stretchr/testify/assert"
)

func TestBuildHarfbuzzString(t *testing.T) {
	tests := []struct {
		subset   string
		expected string
	}{
		{
			subset:   "latin",
			expected: "0000-00FF\n0131\n0152-0153\n02BB-02BC\n02C6\n02DA\n02DC\n0304\n0308\n0329\n2000-206F\n20AC\n2122\n2191\n2193\n2212\n2215\nFEFF\nFFFD\n",
		},
		{
			subset:   "latin-ext",
			expected: "0100-02BA\n02BD-02C5\n02C7-02CC\n02CE-02D7\n02DD-02FF\n0304\n0308\n0329\n1D00-1DBF\n1E00-1E9F\n1EF2-1EFF\n2020\n20A0-20AB\n20AD-20C0\n2113\n2C60-2C7F\nA720-A7FF\n",
		},
		{
			subset:   "vietnamese",
			expected: "0102-0103\n0110-0111\n0128-0129\n0168-0169\n01A0-01A1\n01AF-01B0\n0300-0301\n0303-0304\n0308-0309\n0323\n0329\n1EA0-1EF9\n20AB\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.subset, func(t *testing.T) {
			assert.Equal(t, tt.expected, subsetting.BuildHarfbuzzString(tt.subset), "Failed to build %s", tt.subset)
		})
	}
}

func TestBuildCSSString(t *testing.T) {
	tests := []struct {
		subset   string
		expected string
	}{
		{
			subset:   "latin",
			expected: "U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD",
		},
		{
			subset:   "latin-ext",
			expected: "U+0100-02BA, U+02BD-02C5, U+02C7-02CC, U+02CE-02D7, U+02DD-02FF, U+0304, U+0308, U+0329, U+1D00-1DBF, U+1E00-1E9F, U+1EF2-1EFF, U+2020, U+20A0-20AB, U+20AD-20C0, U+2113, U+2C60-2C7F, U+A720-A7FF",
		},
		{
			subset:   "vietnamese",
			expected: "U+0102-0103, U+0110-0111, U+0128-0129, U+0168-0169, U+01A0-01A1, U+01AF-01B0, U+0300-0301, U+0303-0304, U+0308-0309, U+0323, U+0329, U+1EA0-1EF9, U+20AB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.subset, func(t *testing.T) {
			assert.Equal(t, tt.expected, subsetting.BuildCSSString(tt.subset), "Failed to build %s", tt.subset)
		})
	}
}

func TestInvalidSubsetKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid subset key, but code did not panic")
		}
	}()

	subsetting.BuildHarfbuzzString("invalid-key")
}

func TestInvalidSubsetKeyCSS(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid subset key, but code did not panic")
		}
	}()

	subsetting.BuildCSSString("invalid-key")
}
