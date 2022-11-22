package classify

import (
	"regexp"

	"github.com/bearer/curio/pkg/report/detectors"
)

type ValidationState string

var regexpVendoredFilenameMatcher = regexp.MustCompile(`(vendor/)|(migrations/)|(public/)`)

var Valid = ValidationState("valid")
var Invalid = ValidationState("invalid")
var Potential = ValidationState("potential")

type ValidationResult struct {
	State  ValidationState
	Reason string
}

type ClassificationDecision struct {
	State  ValidationState `json:"state" yaml:"state"`
	Reason string          `json:"reason" yaml:"reason"`
}

var potentialDetectors = map[string]struct{}{
	"env_file":    {},
	"yaml_config": {},
}

const IncludedInVendorFolderReason = "included_in_vendor_folder"
const PotentialDetectorReason = "potential_detectors"
const UnsupportedDetectorType = "unsupported_filename"

func IsVendored(filename string) bool {
	return regexpVendoredFilenameMatcher.MatchString(filename)
}

func IsPotentialDetector(detectorType detectors.Type) bool {
	_, ok := potentialDetectors[string(detectorType)]
	return ok
}
