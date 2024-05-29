package file_errors

import (
	"github.com/bearer/bearer/pkg/report/detections"
	fileerrors "github.com/bearer/bearer/pkg/report/output/dataflow/types"
)

type Holder struct {
	errors []fileerrors.Error
}

func New() *Holder {
	return &Holder{
		errors: []fileerrors.Error{},
	}
}

func (holder *Holder) AddFileError(detection detections.FileFailedDetection) {
	holder.errors = append(holder.errors, fileerrors.Error{
		Type:     string(detection.Type),
		Filename: detection.File,
		Error:    detection.Error,
	})
}

func (holder *Holder) AddError(detection detections.ErrorDetection) {
	holder.errors = append(holder.errors, fileerrors.Error{
		Type:     string(detection.Type),
		Filename: detection.File,
		Error:    detection.Message,
	})
}

func (holder *Holder) ToDataFlow() []fileerrors.Error {
	return holder.errors
}
