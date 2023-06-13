package html

import (
	"strings"
	"text/template"

	html "github.com/bearer/bearer/pkg/report/output/html/types"
	privacy "github.com/bearer/bearer/pkg/report/output/privacy"
	security "github.com/bearer/bearer/pkg/report/output/security"
	"github.com/russross/blackfriday"
)

func ReportHTMLWrapper(body *string) (*string, error) {
	htmlContent := &strings.Builder{}

	wrapperContent := html.WrapperHTMLPage{
		Body:  *body,
		Title: "Bearer Report",
	}

	htmlTemplate := `
<!DOCTYPE html>
<html lang="en-US">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width" />
    <title>{{.Title}}</title>
  </head>
  <body>
		{{.Body}}
  </body>
</html>
`
	pageTemplate, err := template.New("pageTemplate").Parse(htmlTemplate)
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

	htmlTemplate := `
	<ul>
		{{range $severity, $results := .}}
		{{range $index, $result := $results}}
			<li>
				<h2>[{{$severity}}] <a href="{{.Rule.DocumentationUrl}}" target="_blank">{{.Rule.Title}}</a></h2>
				<p><strong>CWE IDs:</strong> {{range .Rule.CWEIDs}}{{.}} {{end}}</p>
				<p><strong>Filename:</strong> {{.Filename}}:{{.LineNumber}}</p>

				{{.Rule.Description | markdownToHtml }}
			</li>
		{{end}}
		{{end}}
	</ul>
	`

	findingsTemplate, err := template.New("findingsTemplate").Funcs(template.FuncMap{
		"kebabCase":      KebabCase,
		"markdownToHtml": MarkdownToHtml,
	}).Parse(htmlTemplate)
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

	for dataSubjectName, subjectGroup := range subjectGroups {
		group := html.GroupedDataSubject{
			DataSubjectName: dataSubjectName,
			Subject:         subjectGroup,
		}
		privacyPage.GroupedDataSubject = append(privacyPage.GroupedDataSubject, group)
	}

	thirdPartyGroups := make(map[string][]privacy.ThirdParty)
	for _, thirdParty := range privacyReport.ThirdParty {
		thirdPartyGroups[thirdParty.ThirdParty] = append(thirdPartyGroups[thirdParty.ThirdParty], thirdParty)
	}

	for thirdPartyName, thirdPartyGroup := range thirdPartyGroups {
		group := html.GroupedThirdParty{
			ThirdPartyName: thirdPartyName,
			ThirdParty:     thirdPartyGroup,
		}
		privacyPage.GroupedThirdParty = append(privacyPage.GroupedThirdParty, group)
	}

	htmlTemplate := `
	<h1>Data Subjects</h1>
	{{- range .GroupedDataSubject -}}
		<h2 id="{{.DataSubjectName | kebabCase }}">{{.DataSubjectName }}</h2>
		<table>
			<tr>
				<th>Data Type</th>
				<th>Detection Count</th>
				<th>Critical Risk Finding Count</th>
				<th>High Risk Finding Count</th>
				<th>Medium Risk Finding Count</th>
				<th>Low Risk Finding Count</th>
				<th>Rules Passed Count</th>
			</tr>
		{{- range .Subject -}}
			<tr>
				<td>{{.DataType}}</td>
				<td>{{.DetectionCount}}</td>
				<td>{{.CriticalRiskFindingCount}}</td>
				<td>{{.HighRiskFindingCount}}</td>
				<td>{{.MediumRiskFindingCount}}</td>
				<td>{{.LowRiskFindingCount}}</td>
				<td>{{.RulesPassedCount}}</td>
			</tr>
		{{- end -}}
		</table>
	{{- end -}}
	<h1>Third Parties</h1>
	{{- range .GroupedThirdParty -}}
		<h2>{{.ThirdPartyName}}</h2>
		<table>
			<tr>
				<th>Subject / Data Type</th>
				<th>Critical Risk Finding Count</th>
				<th>High Risk Finding Count</th>
				<th>Medium Risk Finding Count</th>
				<th>Low Risk Finding Count</th>
				<th>Rules Passed Count</th>
			</tr>
			{{- range .ThirdParty -}}
				<tr>
					<td><a href="#{{.DataSubject | kebabCase }}">{{.DataSubject}}</a> -
					({{- range .DataTypes -}}
						{{.}}&nbsp;
					{{- end -}})</td>
					<td>{{.CriticalRiskFindingCount}}</td>
					<td>{{.HighRiskFindingCount}}</td>
					<td>{{.MediumRiskFindingCount}}</td>
					<td>{{.LowRiskFindingCount}}</td>
					<td>{{.RulesPassedCount}}</td>
				</tr>
		{{- end -}}
		</table>
	{{- end -}}
	`
	subjectTemplate, err := template.New("subjectTemplate").Funcs(template.FuncMap{
		"kebabCase": KebabCase,
	}).Parse(htmlTemplate)
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
	return string(html)
}
