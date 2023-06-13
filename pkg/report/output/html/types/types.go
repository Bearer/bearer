package types

import (
	privacy "github.com/bearer/bearer/pkg/report/output/privacy"
)

type GroupedThirdParty struct {
	ThirdPartyName string
	ThirdParty     []privacy.ThirdParty
}
type GroupedDataSubject struct {
	DataSubjectName string
	Subject         []privacy.Subject
}

type PrivacyHTMLBody = struct {
	GroupedDataSubject []GroupedDataSubject
	GroupedThirdParty  []GroupedThirdParty
}

type WrapperHTMLPage = struct {
	Body  string
	Title string
}
