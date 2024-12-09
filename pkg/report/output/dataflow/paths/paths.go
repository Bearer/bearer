package paths

import (
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/operations/operationshelper"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
)

type Holder struct {
	detectors map[string]*detector
}

type detector struct {
	detections []*types.Detection
	name       string
}

func New(isInternal bool) *Holder {
	return &Holder{
		detectors: make(map[string]*detector),
	}
}

func (holder *Holder) AddOperation(detectorType detectors.Type, detection operationshelper.Operation, fullFilename string) {
	var urls []string

	for _, url := range detection.Value.Urls {
		urls = append(urls, url.Url)
	}

	holder.addPath(
		string(detectorType),
		detection.Source.Filename,
		fullFilename,
		*detection.Source.StartLineNumber,
		detection.Value.Type,
		detection.Value.Path,
		urls,
	)
}

func (holder *Holder) addPath(
	detectorName string,
	fileName string,
	fullFilename string,
	lineNumber int,
	httpMethod string,
	path string,
	urls []string,
) {

	if _, exists := holder.detectors[detectorName]; !exists {
		holder.detectors[detectorName] = &detector{
			name:       detectorName,
			detections: []*types.Detection{},
		}
	}

	targetDetector := holder.detectors[detectorName]
	targetDetector.detections = append(targetDetector.detections, &types.Detection{
		FullFilename: fullFilename,
		FullName:     fileName,
		LineNumber:   &lineNumber,
		HttpMethod:   httpMethod,
		Path:         path,
		Urls:         urls,
	})
}

func (holder *Holder) ToDataFlow() []types.Path {
	data := make([]types.Path, 0)
	for _, detector := range holder.detectors {
		data = append(data, types.Path{
			DetectorName: detector.name,
			Detections:   holder.detectors[detector.name].detections,
		})
	}

	return data
}
