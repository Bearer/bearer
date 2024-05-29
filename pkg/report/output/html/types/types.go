package types

import (
	privacytypes "github.com/bearer/bearer/pkg/report/output/privacy/types"
)

type GroupedThirdParty struct {
	ThirdPartyName string
	ThirdParty     []privacytypes.ThirdParty
}
type GroupedDataSubject struct {
	DataSubjectName string
	Subject         []privacytypes.Subject
}

type PrivacyHTMLBody = struct {
	GroupedDataSubject []GroupedDataSubject
	GroupedThirdParty  []GroupedThirdParty
}

type WrapperHTMLPage = struct {
	Body      string
	Title     string
	TimeStamp string
	Style     string
}
