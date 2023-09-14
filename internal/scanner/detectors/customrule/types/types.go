package types

import (
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/variableshape"
)

type Data struct {
	Pattern   string
	Datatypes []*detectortypes.Detection
	Variables variableshape.Values
}
