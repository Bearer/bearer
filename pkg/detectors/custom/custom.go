package custom

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/detectors/custom/config"
	rubydetector "github.com/bearer/curio/pkg/detectors/ruby/custom_detector"
	rubydatatype "github.com/bearer/curio/pkg/detectors/ruby/datatype"
	sqldetector "github.com/bearer/curio/pkg/detectors/sql/custom_detector"
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	language "github.com/bearer/curio/pkg/parser/custom"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/report/detections"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"
	pluralize "github.com/gertd/go-pluralize"
	"golang.org/x/exp/slices"

	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/parser/sitter/sql"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/smacker/go-tree-sitter/ruby"
)

var insecureUrlPattern = regexp.MustCompile(`^http[^s]`)
var createTableRegexp = regexp.MustCompile(`(?i)(create table)`)
var sqlLanguage = sql.GetLanguage()

type Detector struct {
	idGenerator        nodeid.Generator
	paramIdGenerator   nodeid.Generator
	rulesGroupedByLang map[string][]config.CompiledRule
	pluralize          *pluralize.Client

	Ruby language.Detector
	Sql  language.Detector
}

func New(idGenerator nodeid.Generator) types.Detector {
	detector := &Detector{
		idGenerator:      idGenerator,
		paramIdGenerator: &nodeid.IntGenerator{Counter: 1},
	}

	detector.Ruby = &rubydetector.Detector{}
	detector.Sql = &sqldetector.Detector{}
	detector.pluralize = pluralize.NewClient()

	return detector
}

func (detector *Detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *Detector) CompileRules(rulesConfig map[string]settings.Rule) error {
	detector.rulesGroupedByLang = make(map[string][]config.CompiledRule)

	for ruleName, rule := range rulesConfig {
		if rule.Disabled {
			continue
		}
		for _, lang := range rule.Languages {
			for _, rulePattern := range rule.Patterns {
				compiledRule, err := detector.compileRule(rulePattern, lang, detector.paramIdGenerator)
				if err != nil {
					return err
				}

				compiledRule.Query = parser.QueryMustCompile(getLanguage(lang), compiledRule.Tree)
				compiledRule.Language = lang

				compiledRule.RuleName = ruleName
				compiledRule.Metavars = rule.Metavars
				compiledRule.ParamParenting = rule.ParamParenting
				if rule.ParamParenting {
					compiledRule.VariableReconciliation = false
				} else {
					compiledRule.VariableReconciliation = true
				}
				compiledRule.RootLowercase = rule.RootLowercase
				compiledRule.RootSingularize = rule.RootSingularize
				compiledRule.DetectPresence = rule.DetectPresence
				compiledRule.Pattern = rulePattern.Pattern
				compiledRule.Filters = rulePattern.Filters
				compiledRule.OmitParent = rule.OmitParent

				detector.rulesGroupedByLang[lang] = append(detector.rulesGroupedByLang[lang], compiledRule)
			}
		}
	}

	for _, rules := range detector.rulesGroupedByLang {
		sort.Slice(rules, func(i, j int) bool {
			return strings.Compare(rules[i].RuleName, rules[j].RuleName) == -1
		})
	}

	return nil
}

func (detector *Detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	for lang, rules := range detector.rulesGroupedByLang {
		if !languageMatchesFile(file, lang) {
			continue
		}

		switch lang {
		case "sql":
			f, err := os.Open(file.Path.AbsolutePath)
			if err != nil {
				return false, err
			}
			fileBytes, err := io.ReadAll(f)
			if err != nil {
				return false, err
			}
			defer f.Close()
			// our sql tree sitter parser tends to error sometimes mid file causing us to partially parse file
			// with this hack we increase our parsing percentage
			chunks := createTableRegexp.Split(string(fileBytes), -1)

			lineOffset := 0
			for i, chunk := range chunks {
				chunkBytes := []byte(chunk)
				if i != 0 {
					chunkBytes = []byte("CREATE TABLE" + chunk)
				}

				tree, err := parser.ParseBytes(file, file.Path, chunkBytes, sqlLanguage, lineOffset)
				if err != nil {
					return false, err
				}
				defer tree.Close()

				for _, rule := range rules {
					err := detector.executeRule(rule, tree, report, detector.idGenerator, nil)
					if err != nil {
						return false, err
					}
				}

				lineOffset = lineOffset + strings.Count(chunk, "\n")
			}

			return true, nil
		case "ruby":
			sitterLang := getLanguage(lang)
			tree, err := parser.ParseFile(file, file.Path, sitterLang)
			if err != nil {
				return false, err
			}

			langDetector, err := detector.forLanguage(lang)
			if err != nil {
				return false, err
			}

			if err := langDetector.Annotate(tree); err != nil {
				return false, err
			}

			var variableReconciliation *parserdatatype.ReconciliationRequest

			for _, rule := range rules {
				if rule.VariableReconciliation && variableReconciliation == nil {
					variableReconciliation, err = detector.buildReconciliationRequest(rule.Language, tree)
					if err != nil {
						return false, err
					}
				}

				err := detector.executeRule(rule, tree, report, detector.idGenerator, variableReconciliation)
				if err != nil {
					return false, err
				}
			}

			tree.Close()
		}
	}

	return false, nil
}

func (detector *Detector) forLanguage(lang string) (language.Detector, error) {
	switch lang {
	case "ruby":
		return detector.Ruby, nil
	case "sql":
		return detector.Sql, nil
	default:
		return nil, fmt.Errorf("unsupported language %s", lang)
	}
}

func (detector *Detector) compileRule(
	rulePattern settings.RulePattern,
	lang string,
	idGenerator nodeid.Generator,
) (config.CompiledRule, error) {
	langDetector, err := detector.forLanguage(lang)
	if err != nil {
		return config.CompiledRule{}, err
	}

	rule, err := langDetector.CompilePattern(rulePattern, detector.paramIdGenerator)
	if err != nil {
		return config.CompiledRule{}, err
	}

	return rule, validateRule(rule)
}

func validateRule(rule config.CompiledRule) error {
	for _, filter := range rule.Filters {
		if param := rule.GetParamByPatternName(filter.Variable); param == nil {
			return fmt.Errorf("undefined variable '%s' in filter for custom rule '%s'", filter.Variable, rule.RuleName)
		}
	}

	return nil
}

func languageMatchesFile(file *file.FileInfo, ruleLanguage string) bool {
	if file.Language == "Ruby" && ruleLanguage == "ruby" {
		return true
	}
	if ruleLanguage == "sql" {
		if file.Language != "SQL" &&
			// postgress
			file.Language != "PLpgSQL" && file.Language != "PLSQL" && file.Language != "SQLPL" &&
			// microsoft sql
			file.Language != "TSQL" {
			return false
		}
		return true
	}
	return false
}

func (detector *Detector) executeRule(rule config.CompiledRule, tree *parser.Tree, report report.Report, idGenerator nodeid.Generator, variableReconciliation *parserdatatype.ReconciliationRequest) error {
	captures := tree.QueryMustPass(rule.Query)

	filteredCaptures, err := filterCaptures(rule.Params, captures)
	if err != nil {
		return err
	}

	err = detector.extractData(filteredCaptures, rule, report, rule.Language, idGenerator, variableReconciliation)
	if err != nil {
		return err
	}
	return nil
}

func (detector *Detector) extractData(captures []parser.Captures, rule config.CompiledRule, report report.Report, lang string, idGenerator nodeid.Generator, variableReconciliation *parserdatatype.ReconciliationRequest) error {
	for _, capture := range captures {
		forExport := make(map[parser.NodeID]*schemadatatype.DataType)
		var parent schemadatatype.DataTypable

		if filtersMatch := shouldIgnoreCaptures(capture, rule); filtersMatch {
			continue
		}

		if filtersMatch := shouldMatchCaptures(capture, rule); !filtersMatch {
			continue
		}

		reject := false

		for _, param := range rule.Params {
			var paramTypes map[parser.NodeID]*schemadatatype.DataType
			var err error

			if param.MatchInsecureUrl {
				matchNode := capture[param.BuildFullName()]
				value := matchNode.Value()

				if value != nil && !insecureUrlPattern.MatchString(value.ToString()) {
					reject = true
					break
				}
			}

			if rule.DetectPresence {
				continue
			}

			if param.ArgumentsExtract || param.ClassNameExtract {
				paramTypes, err = detector.extractArguments(lang, capture[param.BuildFullName()], idGenerator, variableReconciliation)
				if err != nil {
					return err
				}

				if rule.ParamParenting {
					// join it as children to master param and export it
					if len(forExport) == 1 {
						for _, datatype := range paramTypes {
							parent.GetProperties()[datatype.GetName()] = datatype
						}

						continue
					}

					// set parent
					if len(paramTypes) == 1 {
						for dataTypeID, datatype := range paramTypes {
							parent, err = parserdatatype.DeepestSingleChild(datatype)
							if err != nil {
								return err
							}

							forExport[dataTypeID] = datatype
							break
						}
						continue
					}
				}

				// set all for export
				for datatypeID, datatype := range paramTypes {
					forExport[datatypeID] = datatype
				}
			}

			if param.StringExtract {
				for _, metavar := range rule.Metavars {
					if metavar.Input == param.PatternName {
						matchNode := capture[param.BuildFullName()]
						content := matchNode.Content()
						searchRegexp := regexp.MustCompile(metavar.Regex)

						matches := searchRegexp.FindAllStringSubmatch(content, -1)

						for _, subgroupMatch := range matches {
							match := subgroupMatch[metavar.Output]

							matchType := &schemadatatype.DataType{
								Node:       matchNode,
								Name:       string(match),
								Type:       schema.SimpleTypeUnknown,
								Properties: make(map[string]schemadatatype.DataTypable),
							}
							matchNodeID := matchNode.ID()

							if rule.ParamParenting {
								parent.GetProperties()[matchType.Name] = matchType
								continue
							}

							forExport[matchNodeID] = matchType
						}
					}
				}
			}
		}

		if reject {
			continue
		}

		if rule.DetectPresence {
			content := capture["rule"].Source(false)
			content.Text = &rule.Pattern

			var parent *schema.Parent
			var parentSource source.Source
			if !rule.OmitParent {
				parentSource = capture["rule"].Source(true)
				parent = &schema.Parent{
					LineNumber: *parentSource.LineNumber,
					Content:    *parentSource.Text,
				}
			} else {
				parentSource = capture["rule"].Source(false)
				parent = &schema.Parent{
					LineNumber: *parentSource.LineNumber,
				}
			}

			report.AddDetection(detections.TypeCustomRisk, detectors.Type(rule.RuleName), content, parent)

			continue
		}

		detector.applyDatatypeTransformations(rule, forExport)

		report.AddDataType(
			detections.TypeCustom,
			detectors.Type(rule.RuleName),
			idGenerator,
			forExport,
			capture["rule"],
		)
	}

	return nil
}

func shouldMatchCaptures(captures parser.Captures, rule config.CompiledRule) bool {
	hasMatchViolationParam := false
	for _, filter := range rule.Filters {
		param := rule.GetParamByPatternName(filter.Variable)
		matchNode := captures[param.BuildFullName()]
		content := matchNode.Content()

		if !filter.MatchViolation {
			continue
		}
		hasMatchViolationParam = true

		if filter.Minimum != nil {
			contentCast, err := strconv.Atoi(content)
			if err != nil {
				return false
			}

			if *filter.Minimum > contentCast {
				return true
			}
		}

		if filter.Maximum != nil {
			contentCast, err := strconv.Atoi(content)
			if err != nil {
				return false
			}

			if *filter.Maximum < contentCast {
				return true
			}
		}
	}

	if hasMatchViolationParam {
		return false
	} else {
		return true
	}
}

func shouldIgnoreCaptures(captures parser.Captures, rule config.CompiledRule) bool {
	for _, filter := range rule.Filters {
		param := rule.GetParamByPatternName(filter.Variable)
		matchNode := captures[param.BuildFullName()]
		content := matchNode.Content()

		if len(filter.Values) > 0 {
			if !slices.Contains(filter.Values, content) {
				return true
			}
		}
	}

	return false
}

func (detector *Detector) applyDatatypeTransformations(rule config.CompiledRule, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	for _, datatype := range datatypes {
		if rule.RootSingularize {
			datatype.Name = detector.pluralize.Singular(datatype.Name)
		}

		if rule.RootLowercase {
			datatype.Name = strings.ToLower(datatype.Name)
		}
	}
}

func filterCaptures(params []config.Param, captures []parser.Captures) (filtered []parser.Captures, err error) {
	for _, capture := range captures {

		// filter captures not matching rules
		shouldIgnore := false
		for _, param := range params {
			if param.StringMatch != "" {
				if capture[param.BuildFullName()].Content() != param.StringMatch {
					shouldIgnore = true
					break
				}
			}

			if param.RegexMatch != "" {
				// todo regex compilation must happen as prestep
				regex, err := regexp.Compile(param.RegexMatch)
				if err != nil {
					return nil, err
				}

				paramContent := capture[param.BuildFullName()].Content()

				if !regex.Match([]byte(paramContent)) {
					shouldIgnore = true
					break
				}

			}
		}

		if shouldIgnore {
			continue
		}

		filtered = append(filtered, capture)
	}

	return filtered, err
}

func (detector *Detector) extractArguments(language string, node *parser.Node, idGenerator nodeid.Generator, variableReconciliation *parserdatatype.ReconciliationRequest) (map[parser.NodeID]*schemadatatype.DataType, error) {
	switch language {
	case "ruby":
		return detector.Ruby.ExtractArguments(node, idGenerator, variableReconciliation)
	case "sql":
		return detector.Sql.ExtractArguments(node, idGenerator, variableReconciliation)
	}
	return nil, errors.New("unsupported language")
}

func (detector *Detector) buildReconciliationRequest(language string, tree *parser.Tree) (*parserdatatype.ReconciliationRequest, error) {
	switch language {
	case "ruby":
		allDatatypes := rubydatatype.Discover(tree.RootNode(), detector.idGenerator)
		scopedDatatypes := parserdatatype.ScopeDatatypes(allDatatypes, detector.idGenerator, rubydatatype.ScopeTerminators)
		return &parserdatatype.ReconciliationRequest{
			ScopedDatatypes:  scopedDatatypes,
			ScopeTerminators: rubydatatype.ScopeTerminators,
		}, nil
	case "sql":
		return nil, nil
	}
	return nil, errors.New("unsupported language")
}

func getLanguage(input string) *sitter.Language {
	switch input {
	case "ruby":
		return ruby.GetLanguage()
	case "sql":
		return sql.GetLanguage()
	}
	return nil
}
