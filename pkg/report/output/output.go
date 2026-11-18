package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-enry/go-enry/v2"
	"github.com/google/uuid"
	"github.com/hhatto/gocloc"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/saas"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/report/output/types"
	globaltypes "github.com/bearer/bearer/pkg/types"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	outputhandler "github.com/bearer/bearer/pkg/util/output"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func GetData(
	report globaltypes.Report,
	config settings.Config,
	gitContext *gitrepository.Context,
	baseBranchFindings *basebranchfindings.Findings,
) (*types.ReportData, error) {
	data := &types.ReportData{}

	// add languages
	languages := make(map[string]int32)
	languageFiles := make(map[string]int32)
	if report.Inputgocloc != nil {
		for _, language := range report.Inputgocloc.Languages {
			languages[language.Name] = language.Code
			languageFiles[language.Name] = int32(len(language.Files))
		}
	}
	data.FoundLanguages = languages
	data.LanguageFiles = languageFiles

	if config.Report.IncludeStats {
		var err error
		data.LanguageStats, err = computeLanguageStats(config.Scan.Target, report.Inputgocloc, config.ReportInputs)
		if err != nil {
			return nil, err
		}
	}

	// add detectors
	if !config.Scan.Quiet && report.HasFiles {
		outputhandler.StdErrLog("Running Detectors")
	}
	err := readReportDetections(data, report.Path)
	if err != nil {
		return data, err
	}
	if config.Report.Report == flag.ReportDetectors {
		// Validate common metadata without adding IDs or changing detector output.
		err = dataflow.AddCoverageDiagnostics(data)
		return data, err
	}

	// add dataflow to data
	if err = GetDataflow(data, report, config, true); err != nil {
		return data, err
	}

	// add report-specific items
	switch config.Report.Report {
	case flag.ReportDataFlow:
		return data, err
	case flag.ReportSecurity:
		err = security.AddReportData(data, config, baseBranchFindings, report.HasFiles)
	case flag.ReportSaaS:
		if err = security.AddReportData(data, config, baseBranchFindings, report.HasFiles); err != nil {
			return nil, err
		}
		err = saas.GetReport(data, config, gitContext, false)
	case flag.ReportPrivacy:
		err = privacy.AddReportData(data, config)
	case flag.ReportStats:
		err = stats.AddReportData(data, report.Inputgocloc, config)
	default:
		return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	return data, err
}

func GetDataflow(
	reportData *types.ReportData,
	report globaltypes.Report,
	config settings.Config,
	isInternal bool,
) error {
	if reportData.Detectors == nil {
		if err := readReportDetections(reportData, report.Path); err != nil {
			return err
		}
	}
	for _, detection := range reportData.Detectors {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}
	return dataflow.AddReportData(reportData, config, isInternal, report.HasFiles)
}

const languageSampleLimit = 16 * 1024

func computeLanguageStats(root string, result *gocloc.Result, scope *inputfile.InputScope) ([]types.LanguageStats, error) {
	if result == nil || len(result.Files) == 0 {
		return nil, nil
	}

	if scope == nil {
		return nil, errors.New("statistics input scope unavailable")
	}
	statsByLanguage := make(map[string]types.LanguageStats)

	absRoot := root
	if absRoot != "" {
		if abs, err := filepath.Abs(root); err == nil {
			absRoot = abs
		}
	}

	var totalBytes int64

	for path, file := range result.Files {
		fullPath := path
		if !filepath.IsAbs(fullPath) {
			if absRoot != "" {
				fullPath = filepath.Join(absRoot, path)
			} else if abs, err := filepath.Abs(path); err == nil {
				fullPath = abs
			}
		}

		relPath := path
		if absRoot != "" {
			if rel, err := filepath.Rel(absRoot, fullPath); err == nil {
				relPath = rel
			}
		}
		relPath = filepath.ToSlash(relPath)

		if shouldSkipLanguagePath(relPath) {
			continue
		}

		sample, fileInfo, err := scope.ReportSample(filepath.FromSlash(relPath), languageSampleLimit)
		if err != nil {
			return nil, fmt.Errorf("statistics sample: %w", err)
		}
		language := detectLanguage(fullPath, file.Lang, sample)
		if language == "" {
			continue
		}

		languageType := enry.GetLanguageType(language)
		if languageType != enry.Programming && languageType != enry.Markup {
			continue
		}

		aggregate := statsByLanguage[language]
		aggregate.Language = language
		aggregate.Lines += file.Code
		aggregate.Files++
		aggregate.Bytes += fileInfo.Size()
		statsByLanguage[language] = aggregate
		totalBytes += fileInfo.Size()
	}

	if len(statsByLanguage) == 0 {
		return nil, nil
	}

	languages := make([]string, 0, len(statsByLanguage))
	for language := range statsByLanguage {
		languages = append(languages, language)
	}
	sort.Strings(languages)

	stats := make([]types.LanguageStats, 0, len(languages))
	for _, language := range languages {
		entry := statsByLanguage[language]
		if totalBytes > 0 {
			entry.Percent = math.Round(float64(entry.Bytes)/float64(totalBytes)*10000) / 100
		}
		stats = append(stats, entry)
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Bytes == stats[j].Bytes {
			if stats[i].Lines == stats[j].Lines {
				return stats[i].Language < stats[j].Language
			}
			return stats[i].Lines > stats[j].Lines
		}
		return stats[i].Bytes > stats[j].Bytes
	})

	return stats, nil
}

func shouldSkipLanguagePath(path string) bool {
	normalized := filepath.ToSlash(path)
	if normalized == ".git" || strings.HasPrefix(normalized, ".git/") {
		return true
	}
	return enry.IsVendor(normalized) ||
		enry.IsDotFile(normalized) ||
		enry.IsDocumentation(normalized) ||
		enry.IsGenerated(normalized, nil)
}

func detectLanguage(path string, fallback string, sample []byte) string {
	language := enry.GetLanguage(filepath.Base(path), sample)
	if language == "" || language == enry.OtherLanguage {
		language = fallback
	}
	if language == "" {
		return ""
	}

	if canonical, ok := enry.GetLanguageByAlias(language); ok {
		language = canonical
	}

	if group := enry.GetLanguageGroup(language); group != "" {
		language = group
	}

	return language
}

func FormatOutput(
	reportData *types.ReportData,
	config settings.Config,
	engine engine.Engine,
	goclocResult *gocloc.Result,
	startTime time.Time,
	endTime time.Time,
) (string, error) {
	var formatter types.GenericFormatter
	switch config.Report.Report {
	case flag.ReportDetectors:
		formatter = detectors.NewFormatter(reportData, config)
	case flag.ReportDataFlow:
		formatter = dataflow.NewFormatter(reportData, config)
	case flag.ReportSecurity:
		formatter = security.NewFormatter(reportData, config, engine, goclocResult, startTime, endTime)
	case flag.ReportPrivacy:
		formatter = privacy.NewFormatter(reportData, config)
	case flag.ReportSaaS:
		formatter = saas.NewFormatter(reportData, config)
	case flag.ReportStats:
		formatter = stats.NewFormatter(reportData, config)
	default:
		return "", fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	formatStr, err := formatter.Format(config.Report.Format)
	if err != nil {
		return formatStr, err
	}
	if formatStr == "" {
		return "", fmt.Errorf(`--report flag "%s" does not support --format flag "%s"`, config.Report.Report, config.Report.Format)
	}

	return formatStr, err
}

// Read all records and surface reader errors rather than silently truncating at
// the old line scanner's token limit. Detection payloads remain unchanged.
func readReportDetections(data *types.ReportData, filename string) (retErr error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, f.Close()) }()
	decoder := json.NewDecoder(f)
	decoder.UseNumber()
	data.Detectors = []interface{}{}
	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		var item map[string]interface{}
		recordDecoder := json.NewDecoder(bytes.NewReader(raw))
		recordDecoder.UseNumber()
		if err := recordDecoder.Decode(&item); err != nil {
			return err
		}
		if item == nil {
			return fmt.Errorf("null report record")
		}
		if item["type"] == inputfile.ContextReceiptType {
			if _, err := inputfile.DecodeContextReceipt(raw); err != nil {
				return err
			}
			continue // Internal handoff is never emitted as a finding or current coverage.
		}
		data.Detectors = append(data.Detectors, preserveReportNumbers(item))
	}
}

// Preserve opaque JSON number tokens in both existing detector serializers.
// json.Number alone marshals as a YAML string; explicit scalar tags retain the
// numeric type without a float64 conversion or a payload-specific schema.
type reportNumber string

func (n reportNumber) MarshalJSON() ([]byte, error) { return []byte(n), nil }
func (n reportNumber) MarshalYAML() (interface{}, error) {
	tag := "!!int"
	if strings.ContainsAny(string(n), ".eE") {
		tag = "!!float"
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: string(n)}, nil
}

func preserveReportNumbers(value any) any {
	switch v := value.(type) {
	case json.Number:
		return reportNumber(v.String())
	case map[string]any:
		for key, item := range v {
			v[key] = preserveReportNumbers(item)
		}
	case []any:
		for i, item := range v {
			v[i] = preserveReportNumbers(item)
		}
	}
	return value
}
