package detections

import (
	"time"

	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/source"
)

type DetectionType string

var TypeDependency DetectionType = "dependency"
var TypeDependencyClassified DetectionType = "dependency_classified"
var TypeInterface DetectionType = "interface"
var TypeInterfaceClassified DetectionType = "interface_classified"
var TypeSchema DetectionType = "schema"
var TypeCreateView DetectionType = "create_view"
var TypeOperation DetectionType = "operation"
var TypeFramework DetectionType = "framework"
var TypeFiller DetectionType = "filler"
var TypeError DetectionType = "error"
var TypeFileFailed DetectionType = "file_error"
var TypeSecretleak DetectionType = "secret_leak"
var TypeCustom DetectionType = "custom"

type ReportDetection interface {
	AddDetection(detectionType DetectionType, detectorType detectors.Type, source source.Source, value interface{})
}

// broker writes those for files that scanner fails to proccess
type FileFailedDetection struct {
	Type     DetectionType `json:"type"`
	File     string        `json:"file"`
	FileSize int           `json:"file_size"`
	Timeout  time.Duration `json:"timeout_duration"`
	Error    string        `json:"error"`
}

type ErrorDetection struct {
	Type    DetectionType `json:"type"`
	Message string        `json:"message"`
}

type FrameworkDetection struct {
	Type          DetectionType   `json:"type"`
	DetectorType  detectors.Type  `json:"detector_type"`
	FrameworkType frameworks.Type `json:"framework_detection_type"`
	CommitSHA     string          `json:"commit_sha"`
	Source        source.Source   `json:"source"`
	Value         interface{}     `json:"value"`
}

type Detection struct {
	Type         DetectionType  `json:"type"`
	DetectorType detectors.Type `json:"detector_type"`
	CommitSHA    string         `json:"commit_sha"`
	Source       source.Source  `json:"source"`
	Value        interface{}    `json:"value"`
}
