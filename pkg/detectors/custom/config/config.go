package config

import "github.com/bearer/curio/pkg/commands/process/settings"

type CompiledRule struct {
	RuleName       string
	Tree           string
	Params         []Param
	Metavars       map[string]settings.MetaVar
	ParamParenting bool
	Singularize    bool
	Lowercase      bool
	Languages      []string
}

type Param struct {
	PatternName      string `yaml:"syntax_name"` // name in pattern eg: $SCRIPT
	ParamName        string `yaml:"param_name"`  // name of param eg: var1
	StringMatch      string `yaml:"string_match"`
	RegexMatch       string `yaml:"regex_match"`
	StringExtract    bool   `yaml:"string_extract"`
	ArgumentsExtract bool   `yaml:"arguments_extract"`
	ClassNameExtract bool   `yaml:"class_name_extract"`
}

func (param *Param) BuildFullName() string {
	return "param_" + param.ParamName
}
