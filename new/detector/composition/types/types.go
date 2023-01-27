package types

import (
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/source"

	reportdetections "github.com/bearer/curio/pkg/report/detections"
)

type Detection struct {
	CustomDetector detectors.Type
	DetectionType  reportdetections.DetectionType
	Source         source.Source
	Value          interface{}
}
