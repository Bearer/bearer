package classify

import (
	"fmt"

	"github.com/bearer/bearer/pkg/report/detectors"
)

var objectStopWords = map[string]struct{}{
	"this":        {},
	"props":       {},
	"prop types":  {},
	"exports":     {},
	"export":      {},
	"env":         {},
	"argv":        {},
	"arguments":   {},
	"errors":      {},
	"args":        {},
	"state":       {},
	"filter":      {},
	"memberships": {},
}

var propertyStopWords = map[string]struct{}{
	"on click":      {},
	"disable click": {},
}

var databaseDetectorTypes = map[string]struct{}{
	"sql_lang_create_table": {},
	"rails":                 {},
	"schema_rb":             {},
}

var expectedIdentifierDataTypeIds = map[string]struct{}{
	"132": {}, // Unique Identifier
	"13":  {}, // Device Identifier
}

func IsDatabase(detectorType detectors.Type) bool {
	_, ok := databaseDetectorTypes[string(detectorType)]
	return ok
}

func IsJSDetection(detectorType detectors.Type) bool {
	return detectorType == detectors.DetectorJavascript || detectorType == detectors.DetectorTypescript
}

func ObjectStopWordDetected(name string) bool {
	_, ok := objectStopWords[name]
	return ok
}

func PropertyStopWordDetected(name string) bool {
	_, ok := propertyStopWords[name]
	return ok
}

func IsExpectedIdentifierDataTypeId(id int) bool {
	_, ok := expectedIdentifierDataTypeIds[fmt.Sprint(id)]
	return ok
}
