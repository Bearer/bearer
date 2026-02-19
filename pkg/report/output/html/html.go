package html

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	term "github.com/buildkite/terminal"
	"github.com/russross/blackfriday/v2"

	html "github.com/bearer/bearer/pkg/report/output/html/types"
	privacytypes "github.com/bearer/bearer/pkg/report/output/privacy/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	"github.com/bearer/bearer/pkg/util/maputil"
)

//go:embed security.tmpl
var securityTemplate string

//go:embed privacy.tmpl
var privacyTemplate string

//go:embed wrapper.tmpl
var wrapperTemplate string

//go:embed styles.css
var siteCss string

func ReportHTMLWrapper(title string, body *string) (string, error) {
	htmlContent := &strings.Builder{}

	t := time.Now()
	timeLayout := "January 2 2006, 15:04:05 pm (MST-0700)"

	wrapperContent := html.WrapperHTMLPage{
		Body:      *body,
		Title:     title,
		TimeStamp: t.Format(timeLayout),
		Style:     strings.Trim(siteCss, ""),
	}
	pageTemplate, err := template.New("pageTemplate").Parse(wrapperTemplate)
	if err != nil {
		return "", err
	}

	err = pageTemplate.Execute(htmlContent, wrapperContent)

	if err != nil {
		return "", err
	}

	return htmlContent.String(), nil
}

func ReportSecurityHTML(detections map[string][]securitytypes.Finding) (*string, error) {
	htmlContent := &strings.Builder{}

	findingsTemplate, err := template.New("findingsTemplate").Funcs(template.FuncMap{
		"kebabCase":      kebabCase,
		"markdownToHtml": markdownToHtml,
		"joinCwe":        joinCwe,
		"count":          countItems,
		"displayExtract": displayExtract,
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

func ReportPrivacyHTML(privacyReport *privacytypes.Report) (*string, error) {
	htmlContent := &strings.Builder{}

	privacyPage := html.PrivacyHTMLBody{
		GroupedDataSubject: make([]html.GroupedDataSubject, 0),
		GroupedThirdParty:  make([]html.GroupedThirdParty, 0),
	}

	subjectGroups := make(map[string][]privacytypes.Subject)
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

	thirdPartyGroups := make(map[string][]privacytypes.ThirdParty)
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
		"kebabCase": kebabCase,
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

func kebabCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func markdownToHtml(s string) string {
	html := blackfriday.Run([]byte(s))
	return strings.ReplaceAll(string(html), "<h2", "<h4")
}

func joinCwe(data []string) string {
	var out = []string{}
	for _, cwe := range data {
		out = append(out, "CWE "+cwe)
	}
	return strings.Join(out, ", ")
}

func countItems(arr interface{}) string {
	switch v := arr.(type) {
	case []securitytypes.Finding:
		return fmt.Sprint(len(v))
	default:
		return "0"
	}
}

func displayExtract(finding securitytypes.Finding) string {
	terminalOutput := finding.HighlightCodeExtract()
	return string(term.Render([]byte(terminalOutput)))
}
