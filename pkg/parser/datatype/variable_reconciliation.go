package datatype

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

func VariableReconciliation(singleArgumentDatatypes map[parser.NodeID]*datatype.DataType, allDataTypes map[parser.NodeID]*datatype.DataType, scopeTerminators []string) {
	for _, argumentDatatype := range singleArgumentDatatypes {
		currentNode := argumentDatatype.Node
		for {
			isTerminating := false

			for _, terminator := range scopeTerminators {
				if currentNode.Type() == terminator {
					isTerminating = true
					break
				}
			}

			if !isTerminating {
				continue
			}

			for globalNodeID, globalDatatype := range allDataTypes {
				// not in the same scope
				if globalNodeID != currentNode.ID() {
					continue
				}
				// in the same scope but it doesn't interest us
				if globalDatatype.Name != argumentDatatype.Name {
					continue
				}

				// merge properties of argumentDatatype and globalDatatype
				MergeDatatypesByPropertyNames(argumentDatatype, globalDatatype)
			}

			currentNode = currentNode.Parent()
			if currentNode == nil {
				break
			}
		}
	}
}
