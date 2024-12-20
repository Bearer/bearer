package types

import (
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/variableshape"
)

type Data struct {
	Pattern   string
	Datatypes []*detectortypes.Detection
	Variables variableshape.Values
	Value     string
}
