package config

import (
	"github.com/bearer/bearer/pkg/commands/process/settings"
	sitter "github.com/smacker/go-tree-sitter"
)

type CompiledRule struct {
	RuleName        string
	Tree            string
	Query           *sitter.Query
	Params          []Param
	Metavars        map[string]settings.MetaVar
	Filters         []settings.PatternFilter
	ParamParenting  bool
	RootSingularize bool
	RootLowercase   bool
	Language        string
	DetectPresence  bool
	Pattern         string
	OmitParent      bool
}

func (rule *CompiledRule) GetParamByPatternName(name string) *Param {
	for _, param := range rule.Params {
		if param.PatternName == name {
			return &param
		}
	}

	return nil
}

type Param struct {
	PatternName      string `yaml:"syntax_name"` // name in pattern eg: $SCRIPT
	ParamName        string `yaml:"param_name"`  // name of param eg: var1
	StringMatch      string `yaml:"string_match"`
	RegexMatch       string `yaml:"regex_match"`
	StringExtract    bool   `yaml:"string_extract"`
	ArgumentsExtract bool   `yaml:"arguments_extract"`
	ClassNameExtract bool   `yaml:"class_name_extract"`
	MatchAnything    bool   `yaml:"match_anything"`
	MatchInsecureUrl bool   `yaml:"match_insecure_url"`
}

func (param *Param) BuildFullName() string {
	return "param_" + param.ParamName
}
