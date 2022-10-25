package datatype

import (
	"sort"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

type Scope struct {
	NodeId    parser.NodeID
	Node      *parser.Node
	DataTypes map[string][]*datatype.DataType
}

func (scope *Scope) toSortedDatatypes() [][]*datatype.DataType {
	var sortedDatatypes [][]*datatype.DataType
	for _, datatypes := range scope.DataTypes {
		sortedDatatypes = append(sortedDatatypes, datatypes)
	}

	for i, datatypes := range sortedDatatypes {
		sortedDatatypes[i] = datatype.SortSlice(datatypes)
	}

	sort.Slice(sortedDatatypes, func(i, j int) bool {
		lineNumberA := sortedDatatypes[i][0].Node.Source(false).LineNumber
		lineNumberB := sortedDatatypes[j][0].Node.Source(false).LineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := sortedDatatypes[i][0].Node.Source(false).ColumnNumber
		columnNumberB := sortedDatatypes[j][0].Node.Source(false).ColumnNumber

		return *columnNumberA < *columnNumberB
	})

	return sortedDatatypes
}

func ScopeDatatypes(datatypes map[parser.NodeID]*datatype.DataType, idGenerator nodeid.Generator, termintingTokens []string) {
	scopes := make(map[parser.NodeID]*Scope)

	// iterate trough datatypes
	for _, target := range datatypes {
		scopeNode := target.Node.Parent()

		if scopeNode == nil {
			scopeNode = target.Node
		} else {
			// find scope terminating parent
			for {
				parentType := scopeNode.Type()

				terminatorFound := false
				for _, v := range termintingTokens {
					if v == parentType {
						terminatorFound = true
						break
					}
				}

				if terminatorFound {
					break
				}

				if scopeNode.Parent() == nil {
					break
				}

				scopeNode = scopeNode.Parent()
			}
		}

		// ensure scopedDatatype exist
		_, scopeExists := scopes[scopeNode.ID()]
		if !scopeExists {
			scopes[scopeNode.ID()] = &Scope{
				NodeId:    scopeNode.ID(),
				Node:      scopeNode,
				DataTypes: make(map[string][]*datatype.DataType),
			}
		}

		// append same scope same name datatypes
		scopes[scopeNode.ID()].DataTypes[target.Name] = append(scopes[scopeNode.ID()].DataTypes[target.Name], target)
	}

	var sortedScopes []*Scope
	for _, scope := range scopes {
		sortedScopes = append(sortedScopes, scope)
	}

	SortScopes(sortedScopes)

	for _, scope := range sortedScopes {
		sortedDatatypes := scope.toSortedDatatypes()
		for _, datatypes := range sortedDatatypes {
			UnifyUUID(datatypes, idGenerator)
		}
	}
}

func UnifyUUID(datatypes []*datatype.DataType, idGenerator nodeid.Generator) {
	datatypeID := idGenerator.GenerateId()
	for _, target := range datatypes {
		target.UUID = datatypeID
	}

	propertiesDone := make(map[string]bool)

	for _, target := range datatypes {

		var propertyNames []string
		for propertyName := range target.Properties {
			propertyNames = append(propertyNames, propertyName)
		}

		sort.Strings(propertyNames)

		for _, propertyName := range propertyNames {
			_, alreadyDone := propertiesDone[propertyName]
			if alreadyDone {
				continue
			}

			var datatypesToDo []*datatype.DataType
			// fetch all datatypes that have that property name
			for _, target := range datatypes {
				_, hasProperty := target.Properties[propertyName]
				if hasProperty {
					datatypesToDo = append(datatypesToDo, target.Properties[propertyName])
				}
			}

			UnifyUUID(datatypesToDo, idGenerator)

			propertiesDone[propertyName] = true
		}
	}
}

func SortScopes(input []*Scope) {
	sort.Slice(input, func(i, j int) bool {
		lineNumberA := input[i].Node.Source(false).LineNumber
		lineNumberB := input[j].Node.Source(false).LineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := input[i].Node.Source(false).ColumnNumber
		columnNumberB := input[j].Node.Source(false).ColumnNumber

		return *columnNumberA < *columnNumberB
	})
}
