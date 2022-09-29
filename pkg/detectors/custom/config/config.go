package config

type Config struct {
	Rules map[string]Rule
}

type Rule struct {
	Disabled       bool
	Languages      []string
	Patterns       []string
	ParamParenting bool `yaml:"param_parenting"`
	Metavars       map[string]MetaVar
}

type CompiledRule struct {
	RuleName       string
	Tree           string
	Params         []Param
	Metavars       map[string]MetaVar
	ParamParenting bool
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

type MetaVar struct {
	Input  string
	Output int
	Regex  string
}
