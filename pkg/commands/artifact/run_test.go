package artifact

import (
	"testing"

	"github.com/hhatto/gocloc"
	"github.com/stretchr/testify/assert"
)

func TestFormatLanguagesWithJavascriptAndTypescript(t *testing.T) {
	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{},
		Languages: map[string]*gocloc.Language{
			"Ruby": {
				Name: "Ruby",
			},
			"TypeScript": {
				Name: "TypeScript",
			},
			"JavaScript": {
				Name: "JavaScript",
			},
		},
		MaxPathLength: 0,
	}

	assert.Equal(
		t,
		[]string{"javascript", "ruby"},
		FormatFoundLanguages(dummyGoclocResult.Languages),
	)
}

func TestFormatLanguagesWithoutJavascript(t *testing.T) {
	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{},
		Languages: map[string]*gocloc.Language{
			"Ruby": {
				Name: "Ruby",
			},
			"TypeScript": {
				Name: "TypeScript",
			},
		},
		MaxPathLength: 0,
	}

	assert.Equal(
		t,
		[]string{"javascript", "ruby"},
		FormatFoundLanguages(dummyGoclocResult.Languages),
	)
}

func TestFormatLanguagesWithJavascriptFirst(t *testing.T) {
	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{},
		Languages: map[string]*gocloc.Language{
			"Ruby": {
				Name: "Ruby",
			},
			"JavaScript": {
				Name: "JavaScript",
			},
			"TypeScript": {
				Name: "TypeScript",
			},
		},
		MaxPathLength: 0,
	}

	assert.Equal(
		t,
		[]string{"javascript", "ruby"},
		FormatFoundLanguages(dummyGoclocResult.Languages),
	)
}
