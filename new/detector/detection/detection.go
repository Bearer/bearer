package detection

import "github.com/bearer/bearer/new/language/tree"

type Detection struct {
	DetectorType string
	MatchNode    *tree.Node
	Data         interface{}
}
