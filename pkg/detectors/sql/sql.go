package sql

import (
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/util/file"

	reporttypes "github.com/bearer/curio/pkg/report"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
)

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report reporttypes.Report) (bool, error) {
	// general sql
	if file.Language != "SQL" &&
		// postgress
		file.Language != "PLpgSQL" && file.Language != "PLSQL" && file.Language != "SQLPL" &&
		// microsoft sql
		file.Language != "TSQL" {
		return false, nil
	}

	return true, nil
}

func ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*schemadatatype.DataType, error) {
	return nil, nil
}
