package artifact

import (
	"testing"

	engineimpl "github.com/bearer/bearer/pkg/engine/implementation"
	"github.com/bearer/bearer/pkg/languages"

	"github.com/hhatto/gocloc"
	"github.com/stretchr/testify/assert"
)

func TestGetFoundLanguageIDsWithJavascriptAndTypescript(t *testing.T) {
	engine := engineimpl.New(languages.Default())

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
		GetFoundLanguageIDs(engine, dummyGoclocResult.Languages),
	)
}

func TestGetFoundLanguageIDsWithoutJavascript(t *testing.T) {
	engine := engineimpl.New(languages.Default())

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
		GetFoundLanguageIDs(engine, dummyGoclocResult.Languages),
	)
}

func TestGetFoundLanguageIDsWithJavascriptFirst(t *testing.T) {
	engine := engineimpl.New(languages.Default())

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
		GetFoundLanguageIDs(engine, dummyGoclocResult.Languages),
	)
}
