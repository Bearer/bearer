package html

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	html "github.com/bearer/bearer/pkg/report/output/html/types"
	privacy "github.com/bearer/bearer/pkg/report/output/privacy"
	security "github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/util/maputil"
	term "github.com/buildkite/terminal"
	"github.com/russross/blackfriday"
)

//go:embed security.tmpl
var securityTemplate string

//go:embed privacy.tmpl
var privacyTemplate string

//go:embed wrapper.tmpl
var wrapperTemplate string

//go:embed styles.css
var siteCss string

func ReportHTMLWrapper(title string, body *string) (*string, error) {
	htmlContent := &strings.Builder{}

	t := time.Now()
	timeLayout := "January 2nd 2006, 15:04:05 pm (MST-0700)"

	wrapperContent := html.WrapperHTMLPage{
		Body:      *body,
		Title:     title,
		TimeStamp: t.Format(timeLayout),
		Style:     strings.Trim(siteCss, ""),
	}
	pageTemplate, err := template.New("pageTemplate").Parse(wrapperTemplate)
	if err != nil {
		return nil, err
	}

	err = pageTemplate.Execute(htmlContent, wrapperContent)

	if err != nil {
		return nil, err
	}

	content := htmlContent.String()
	return &content, nil
}

func ReportSecurityHTML(detections *map[string][]security.Result) (*string, error) {
	htmlContent := &strings.Builder{}

	findingsTemplate, err := template.New("findingsTemplate").Funcs(template.FuncMap{
		"kebabCase":      KebabCase,
		"markdownToHtml": MarkdownToHtml,
		"joinCwe":        JoinCwe,
		"count":          CountItems,
		"displayExtract": DisplayExtract,
	}).Parse(securityTemplate)
	if err != nil {
		return nil, err
	}
	err = findingsTemplate.Execute(htmlContent, detections)
	if err != nil {
		return nil, err
	}

	content := htmlContent.String()
	return &content, nil
}

func ReportPrivacyHTML(privacyReport *privacy.Report) (*string, error) {

	htmlContent := &strings.Builder{}

	privacyPage := html.PrivacyHTMLBody{
		GroupedDataSubject: make([]html.GroupedDataSubject, 0),
		GroupedThirdParty:  make([]html.GroupedThirdParty, 0),
	}

	subjectGroups := make(map[string][]privacy.Subject)
	for _, subject := range privacyReport.Subjects {
		subjectGroups[subject.DataSubject] = append(subjectGroups[subject.DataSubject], subject)
	}

	for _, dataSubjectName := range maputil.SortedStringKeys(subjectGroups) {
		group := html.GroupedDataSubject{
			DataSubjectName: dataSubjectName,
			Subject:         subjectGroups[dataSubjectName],
		}
		privacyPage.GroupedDataSubject = append(privacyPage.GroupedDataSubject, group)
	}

	thirdPartyGroups := make(map[string][]privacy.ThirdParty)
	for _, thirdParty := range privacyReport.ThirdParty {
		thirdPartyGroups[thirdParty.ThirdParty] = append(thirdPartyGroups[thirdParty.ThirdParty], thirdParty)
	}

	for _, thirdPartyName := range maputil.SortedStringKeys(thirdPartyGroups) {
		group := html.GroupedThirdParty{
			ThirdPartyName: thirdPartyName,
			ThirdParty:     thirdPartyGroups[thirdPartyName],
		}
		privacyPage.GroupedThirdParty = append(privacyPage.GroupedThirdParty, group)
	}

	subjectTemplate, err := template.New("subjectTemplate").Funcs(template.FuncMap{
		"kebabCase": KebabCase,
	}).Parse(privacyTemplate)
	if err != nil {
		return nil, err
	}

	err = subjectTemplate.Execute(htmlContent, privacyPage)
	if err != nil {
		return nil, err
	}

	content := htmlContent.String()
	return &content, nil
}

func KebabCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func MarkdownToHtml(s string) string {
	html := blackfriday.MarkdownCommon([]byte(s))
	return strings.ReplaceAll(string(html), "<h2", "<h4")
}

func JoinCwe(data []string) string {
	var out = []string{}
	for _, cwe := range data {
		out = append(out, "CWE "+cwe)
	}
	return strings.Join(out, ", ")
}

func CountItems(arr interface{}) string {
	switch v := arr.(type) {
	case []security.Result:
		return fmt.Sprint(len(v))
	default:
		return "0"
	}
}

func DisplayExtract(result security.Result) string {
	terminalOutput := security.HighlightCodeExtract(result)
	return string(term.Render([]byte(terminalOutput)))
}
