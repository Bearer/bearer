package datatype

import (
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
)

func Discover(report report.Report, tree *parser.Tree, idGenerator nodeid.Generator) {
	datatypes := make(map[parser.NodeID]*schemadatatype.DataType)
	helperDatatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	addProperties(tree, helperDatatypes)
	linkProperties(tree, datatypes, helperDatatypes)
	scopeProperties(datatypes, idGenerator)

	addObjects(tree, datatypes)

	parserdatatype.PruneMap(datatypes)

	// remove this
	for nodeID, datatype := range datatypes {
		if datatype.Name == "this" {
			delete(datatypes, nodeID)
		}
	}

	parserdatatype.NewExport(report, detectors.DetectorJavascript, idGenerator, datatypes)
}
