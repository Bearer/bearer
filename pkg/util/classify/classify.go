package classify

import (
	"strings"

	"github.com/bearer/curio/pkg/report/detectors"
)

type ValidationState string

var Valid = ValidationState("valid")
var Invalid = ValidationState("invalid")
var Potential = ValidationState("potential")

type ValidationResult struct {
	State  ValidationState
	Reason string
}

type ClassificationDecision struct {
	State  ValidationState `json:"state"`
	Reason string          `json:"reason"`
}

var potentialDetectors = map[string]struct{}{
	"env_file":    {},
	"yaml_config": {},
}

var cloudDetectors = map[string]struct{}{
	"sql-readonly": {},
	"sql-advanced": {},
}

var supportedDetectorTypes = map[string]struct{}{
	"protobuf": {},
	"proto":    {},
	"graphql":  {},
	"openapi":  {},
	"sql":      {},
	"rails":    {},
}

var extendedLanguageDetectorTypes = map[string]struct{}{
	"ruby":       {},
	"java":       {},
	"csharp":     {},
	"typescript": {},
	"javascript": {},
	"python":     {},
	"go":         {},
	"php":        {},
}

const IncludedInVendorFolderReason = "included_in_vendor_folder"
const PotentialDetectorReason = "potential_detectors"
const UnsupportedDetectorType = "unsupported_filename"

func IsVendored(filename string) bool {
	_, ok := cloudDetectors[filename]
	if ok {
		return false
	}

	return strings.Contains(filename, "vendor/")
}

func IsPotentialDetector(detectorType detectors.Type) bool {
	_, ok := potentialDetectors[string(detectorType)]
	return ok
}

func IsSupportedDetectorType(detectorType detectors.Type) bool {
	_, supportedDetectorType := supportedDetectorTypes[string(detectorType)]
	if supportedDetectorType {
		return true
	}

	_, extendedLanguageDetectorType := extendedLanguageDetectorTypes[string(detectorType)]
	return extendedLanguageDetectorType
}
