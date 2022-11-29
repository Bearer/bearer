package v2yaml

import (
	"github.com/bearer/curio/pkg/detectors/openapi/queries"
	"github.com/bearer/curio/pkg/detectors/openapi/reportadder"
	yamlparser "github.com/bearer/curio/pkg/detectors/openapi/yaml"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	reporttypes "github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/operations/operationshelper"
	"github.com/bearer/curio/pkg/report/schema/schemahelper"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryParameters = parser.QueryMustCompile(yaml.GetLanguage(), `
(_
	(
      block_mapping_pair
        key:
            (flow_node) @helperName
            (#match? @helperName "^name$")
         value:
            (flow_node) @param_name
	)
	(
      block_mapping_pair

        key:
            (flow_node) @helperType
            (#match? @helperType "^type$")
         value:
            (flow_node) @param_type
	)
)
`)

func ProcessFile(idGenerator nodeid.Generator, file *file.FileInfo, report reporttypes.Report) (bool, error) {
	tree, err := parser.ParseFile(file, file.Path, yaml.GetLanguage())
	if err != nil {
		return false, err
	}
	defer tree.Close()

	nodeIDMap := nodeid.New(tree, idGenerator)
	nodeIDMap.Annotate()

	foundValues := make(map[parser.Node]*schemahelper.Schema)

	err = queries.AnnotateV2Paramaters(nodeIDMap, tree, foundValues, queryParameters)
	if err != nil {
		return false, err
	}

	err = yamlparser.AnnotateOperationId(nodeIDMap, tree, foundValues)
	if err != nil {
		return false, err
	}

	err = yamlparser.AnnotateObjects(nodeIDMap, tree, foundValues)
	if err != nil {
		return false, err
	}

	foundPaths := make(map[parser.Node]*operationshelper.Operation)
	err = yamlparser.AnnotatePaths(tree, foundPaths)
	if err != nil {
		return false, err
	}

	reportadder.AddSchema(file, report, foundValues, idGenerator)

	return true, err
}
