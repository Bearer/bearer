package types

import (
	detectortypes "github.com/bearer/curio/new/detector/types"
)

type DetectorInitResult struct {
	Error        error
	Detector     detectortypes.Detector
	DetectorName string
}
