package types

import (
	detectortypes "github.com/bearer/bearer/new/detector/types"
)

type DetectorInitResult struct {
	Error        error
	Detector     detectortypes.Detector
	DetectorName string
	Order        int
}
