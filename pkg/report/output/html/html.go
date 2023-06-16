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
		<style>
		body {
			margin:0;
			background-color: #fff;
			font-family: Source Sans Pro, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
			padding-bottom: 75px;
			font-style: normal;
			font-weight: 400;
			font-size: 14px;
			line-height: 150%;
			/* identical to box height, or 21px */
			letter-spacing: 0.1px;
			/* Neutral/900 */
			color: #272727;
		}
		header {
			background-color: #F1F4FF;
		}
		h1 {
			font-weight: 600;
			font-size: 32px;
			line-height: 110%;
			letter-spacing: 0.3px;
			color: #272727;
			flex: none;
			order: 0;
			flex-grow: 0;
			border-bottom: #EAEAEA 1px solid;
			padding-bottom: 16px;
			margin-bottom: 0px;
		}
		h2 {
			font-weight: 600;
			font-size: 28px;
			line-height: 110%;
			letter-spacing: 0.3px;
			color: #000000;
			flex: none;
			order: 0;
			flex-grow: 0;
			border-bottom: #EAEAEA 1px solid;
			padding-bottom: 16px;
		}
		h3 {
			font-weight: 600;
			font-size: 20px;
			line-height: 140%;
			letter-spacing: 0.2px;
			color: #000000;
			flex: none;
			order: 1;
			flex-grow: 0;
		}
		header{
			height: 75px;
		}
		section {
			margin: auto;
			width:1024px;
		}
		th {
			text-align: left;
			font-style: normal;
			font-weight: 600;
			font-size: 11px;
			line-height: 150%;
			/* identical to box height, or 16px */
			letter-spacing: 0.3px;
			text-transform: uppercase;

			/* Neutral/750 */
			color: #696969;
			padding:11px;


			/* Inside auto layout */
			flex: none;
			order: 0;
			flex-grow: 0;

		}
		td {
			padding:11px;
		}
		table, th, td {
			border: 1px solid #EAEAEA;
			border-collapse: collapse;
		}
		table {
			width: 100%;
		}
		p {
		}

	</style>
  </head>
  <body>
		<header></header>
		<section>{{.Body}}</section>
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
	<h1>Privacy report</h1>
	<p>June 13th 2023, 11:12:23 am (UTC+00:00)</p>
	<h2>Data Subjects</h2>
	{{- range .GroupedDataSubject -}}
		<h3 id="{{.DataSubjectName | kebabCase }}">{{.DataSubjectName }}</h3>
		<table>
			<tr>
				<th>Data Type</th>
				<th>Detection Count</th>
				<th>Critical Risk Findings</th>
				<th>High Risk Findings</th>
				<th>Medium Risk Findings</th>
				<th>Low Risk Findings</th>
				<th>Rules Passed</th>
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
	<h2>Third Parties</h2>
	{{- range .GroupedThirdParty -}}
		<h3>{{.ThirdPartyName}}</h3>
		<table>
			<tr>
				<th>Subject / Data Type</th>
				<th>Critical Risk Findings</th>
				<th>High Risk Findings</th>
				<th>Medium Risk Findings</th>
				<th>Low Risk Findings</th>
				<th>Rules Passed</th>
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
