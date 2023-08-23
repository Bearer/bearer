package datatype

import (
	"github.com/bearer/bearer/internal/parser"
	parserdatatype "github.com/bearer/bearer/internal/parser/datatype"
	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/detectors"
	schemadatatype "github.com/bearer/bearer/internal/report/schema/datatype"
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

	report.AddDataType(detections.TypeSchema, detectors.DetectorJavascript, idGenerator, datatypes, nil)
}
