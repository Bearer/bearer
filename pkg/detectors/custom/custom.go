package custom

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/detectors/custom/config"
	sqldetector "github.com/bearer/bearer/pkg/detectors/sql/custom_detector"
	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	language "github.com/bearer/bearer/pkg/parser/custom"
	parserdatatype "github.com/bearer/bearer/pkg/parser/datatype"
	"github.com/bearer/bearer/pkg/report/detections"
	schemadatatype "github.com/bearer/bearer/pkg/report/schema/datatype"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"

	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/parser/sitter/sql"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
)

var insecureUrlPattern = regexp.MustCompile(`^http[^s]`)
var createTableRegexp = regexp.MustCompile(`(?i)(create table)`)
var sqlLanguage = sql.GetLanguage()

type Detector struct {
	idGenerator        nodeid.Generator
	paramIdGenerator   nodeid.Generator
	rulesGroupedByLang map[string][]config.CompiledRule

	Sql language.Detector
}

func New(idGenerator nodeid.Generator) types.Detector {
	detector := &Detector{
		idGenerator:      idGenerator,
		paramIdGenerator: &nodeid.IntGenerator{Counter: 1},
	}

	detector.Sql = &sqldetector.Detector{}

	return detector
}

func (detector *Detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *Detector) CompileRules(rulesConfig map[string]*settings.Rule) error {
	detector.rulesGroupedByLang = make(map[string][]config.CompiledRule)

	for ruleName, rule := range rulesConfig {
		for _, lang := range rule.Languages {
			for _, rulePattern := range rule.Patterns {
				compiledRule, err := detector.compileRule(rulePattern, lang, detector.paramIdGenerator)
				if err != nil {
					return err
				}

				compiledRule.Query = parser.QueryMustCompile(getLanguage(lang), compiledRule.Tree)
				compiledRule.Language = lang

				compiledRule.RuleName = ruleName
				compiledRule.ParamParenting = rule.ParamParenting
				compiledRule.DetectPresence = rule.DetectPresence
				compiledRule.Pattern = rulePattern.Pattern
				compiledRule.Filters = rulePattern.Filters

				detector.rulesGroupedByLang[lang] = append(detector.rulesGroupedByLang[lang], compiledRule)
			}
		}
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
			f, err := os.Open(file.AbsolutePath)
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
					err := detector.executeRule(rule, tree, report, detector.idGenerator)
					if err != nil {
						return false, err
					}
				}

				lineOffset = lineOffset + strings.Count(chunk, "\n")
			}

			return true, nil
		}
	}

	return false, nil
}

func (detector *Detector) forLanguage(lang string) (language.Detector, error) {
	switch lang {
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

func (detector *Detector) executeRule(rule config.CompiledRule, tree *parser.Tree, report report.Report, idGenerator nodeid.Generator) error {
	captures := tree.QueryMustPass(rule.Query)

	filteredCaptures, err := filterCaptures(rule.Params, captures)
	if err != nil {
		return err
	}

	err = detector.extractData(filteredCaptures, rule, report, rule.Language, idGenerator)
	if err != nil {
		return err
	}
	return nil
}

func (detector *Detector) extractData(captures []parser.Captures, rule config.CompiledRule, report report.Report, lang string, idGenerator nodeid.Generator) error {
	for _, capture := range captures {
		forExport := make(map[parser.NodeID]*schemadatatype.DataType)
		var parent schemadatatype.DataTypable

		if filtersMatch := shouldIgnoreCaptures(capture, rule); filtersMatch {
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
				paramTypes, err = detector.extractArguments(lang, capture[param.BuildFullName()], idGenerator)
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
			var schemaSource *schema.Source
			var source source.Source
			if !rule.OmitParent {
				source = capture["rule"].Source(true)
				schemaSource = &schema.Source{
					StartLineNumber:   *source.StartLineNumber,
					EndLineNumber:     *source.EndLineNumber,
					StartColumnNumber: *source.StartColumnNumber,
					EndColumnNumber:   *source.EndColumnNumber,
				}
			} else {
				source = capture["rule"].Source(false)
				schemaSource = &schema.Source{
					StartLineNumber:   *source.StartLineNumber,
					StartColumnNumber: *source.StartColumnNumber,
					EndLineNumber:     *source.EndLineNumber,
					EndColumnNumber:   *source.EndColumnNumber,
				}
			}

			report.AddDetection(detections.TypeCustomRisk, detectors.Type(rule.RuleName), content, schemaSource)

			continue
		}

		// detector.applyDatatypeTransformations(rule, forExport)

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

				if !regex.MatchString(paramContent) {
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

func (detector *Detector) extractArguments(language string, node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*schemadatatype.DataType, error) {
	switch language {
	case "sql":
		return detector.Sql.ExtractArguments(node, idGenerator)
	}
	return nil, errors.New("unsupported language")
}

func getLanguage(input string) *sitter.Language {
	switch input {
	case "sql":
		return sql.GetLanguage()
	}
	return nil
}
