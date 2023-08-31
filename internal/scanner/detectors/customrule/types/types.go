package types

import (
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
)

type Data struct {
	Pattern       string
	Datatypes     []*detectortypes.Detection
	VariableNodes map[string]*tree.Node
}
