package detections

import (
	"time"

	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/frameworks"
	"github.com/bearer/bearer/pkg/report/source"
)

type DetectionType string

var TypeDependency DetectionType = "dependency"
var TypeDependencyClassified DetectionType = "dependency_classified"
var TypeInterface DetectionType = "interface"
var TypeInterfaceClassified DetectionType = "interface_classified"
var TypeSchema DetectionType = "schema"
var TypeSchemaClassified DetectionType = "schema_classified"
var TypeCreateView DetectionType = "create_view"
var TypeFramework DetectionType = "framework"
var TypeFrameworkClassified DetectionType = "framework_classified"
var TypeFiller DetectionType = "filler"
var TypeError DetectionType = "error"
var TypeFileList DetectionType = "file_list"
var TypeFileFailed DetectionType = "file_error"
var TypeSecretleak DetectionType = "secret_leak"
var TypeCustom DetectionType = "custom"
var TypeCustomClassified DetectionType = "custom_classified"
var TypeCustomRisk DetectionType = "custom_risk"
var TypeExpectedDetection DetectionType = "expected_detection"

type ReportDetection interface {
	AddDetection(detectionType DetectionType, detectorType detectors.Type, source source.Source, value interface{})
}

type FileListDetection struct {
	Type      DetectionType `json:"type" yaml:"type"`
	Filenames []string      `json:"filenames" yaml:"filenames"`
}

// broker writes those for files that scanner fails to proccess
type FileFailedDetection struct {
	Type     DetectionType `json:"type" yaml:"type"`
	File     string        `json:"file" yaml:"file"`
	FileSize int           `json:"file_size" yaml:"file_size"`
	Timeout  time.Duration `json:"timeout_duration" yaml:"timeout_duration"`
	Error    string        `json:"error" yaml:"error"`
}

type ErrorDetection struct {
	Type    DetectionType `json:"type" yaml:"type"`
	Message string        `json:"message" yaml:"message"`
	File    string        `json:"file" yaml:"file"`
}

type FrameworkDetection struct {
	Type          DetectionType   `json:"type" yaml:"type"`
	DetectorType  detectors.Type  `json:"detector_type" yaml:"detector_type"`
	FrameworkType frameworks.Type `json:"framework_detection_type" yaml:"framework_detection_type"`
	CommitSHA     string          `json:"commit_sha,omitempty" yaml:"commit_sha,omitempty"`
	Source        source.Source   `json:"source" yaml:"source"`
	Value         interface{}     `json:"value" yaml:"value"`
}

type Detection struct {
	Type             DetectionType      `json:"type" yaml:"type"`
	DetectorType     detectors.Type     `json:"detector_type" yaml:"detector_type"`
	DetectorLanguage detectors.Language `json:"detector_language,omitempty" yaml:"detector_language,omitempty"`
	CommitSHA        string             `json:"commit_sha,omitempty" yaml:"commit_sha,omitempty"`
	Source           source.Source      `json:"source" yaml:"source"`
	Value            interface{}        `json:"value" yaml:"value"`
}
