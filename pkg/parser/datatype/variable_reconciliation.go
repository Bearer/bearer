package datatype

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

type ReconciliationRequest struct {
	ScopedDatatypes  map[parser.NodeID]*Scope
	ScopeTerminators []string
	Skip             bool
}

func VariableReconciliation(singleArgumentDatatypes map[parser.NodeID]*datatype.DataType, request *ReconciliationRequest) {
	for _, argumentDatatype := range singleArgumentDatatypes {
		currentNode := argumentDatatype.Node
		isFirst := true
		for {
			if isFirst {
				isFirst = false
			} else {
				currentNode = currentNode.Parent()
			}

			if currentNode == nil {
				break
			}

			isTerminating := false

			for _, terminator := range request.ScopeTerminators {
				if currentNode.Type() == terminator {
					isTerminating = true
					break
				}
			}

			if !isTerminating {
				continue
			}

			for scopeID, scope := range request.ScopedDatatypes {

				// not in the same scope
				if scopeID != currentNode.ID() {

					continue
				}

				for datatypeName, datatypeOccurences := range scope.DataTypes {
					// in the same scope but it doesn't interest us because they are different datatypes
					if datatypeName != argumentDatatype.Name {
						continue
					}

					for _, scopedDatatype := range datatypeOccurences {
						// merge properties of argumentDatatype and globalDatatype
						MergeDatatypesByPropertyNames(argumentDatatype, scopedDatatype)
					}
				}
			}
		}
	}
}
