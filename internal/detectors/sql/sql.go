package sql

import (
	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/util/file"

	reporttypes "github.com/bearer/bearer/internal/report"
	schemadatatype "github.com/bearer/bearer/internal/report/schema/datatype"
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
