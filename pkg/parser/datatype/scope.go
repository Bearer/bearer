package datatype

import (
	"sort"

	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/schema/datatype"
)

type Scope struct {
	NodeId    parser.NodeID
	Node      *parser.Node
	DataTypes map[string][]datatype.DataTypable
}

func (scope *Scope) toSortedDatatypes() [][]datatype.DataTypable {
	var sortedDatatypes [][]datatype.DataTypable
	for _, datatypes := range scope.DataTypes {
		sortedDatatypes = append(sortedDatatypes, datatypes)
	}

	for i, datatypes := range sortedDatatypes {
		sortedDatatypes[i] = datatype.SortSlice(datatypes)
	}

	sort.Slice(sortedDatatypes, func(i, j int) bool {
		lineNumberA := sortedDatatypes[i][0].GetNode().Source(false).StartLineNumber
		lineNumberB := sortedDatatypes[j][0].GetNode().Source(false).StartLineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := sortedDatatypes[i][0].GetNode().Source(false).StartColumnNumber
		columnNumberB := sortedDatatypes[j][0].GetNode().Source(false).StartColumnNumber

		return *columnNumberA < *columnNumberB
	})

	return sortedDatatypes
}

func ScopeDatatypes(datatypes map[parser.NodeID]*datatype.DataType, idGenerator nodeid.Generator, termintingTokens []string) map[parser.NodeID]*Scope {
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
				DataTypes: make(map[string][]datatype.DataTypable),
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

	return scopes
}

func UnifyUUID(datatypes []datatype.DataTypable, idGenerator nodeid.Generator) {
	datatypeID := idGenerator.GenerateId()
	for _, target := range datatypes {
		target.SetUUID(datatypeID)
	}

	propertiesDone := make(map[string]bool)

	for _, target := range datatypes {

		var propertyNames []string
		for propertyName := range target.GetProperties() {
			propertyNames = append(propertyNames, propertyName)
		}

		sort.Strings(propertyNames)

		for _, propertyName := range propertyNames {
			_, alreadyDone := propertiesDone[propertyName]
			if alreadyDone {
				continue
			}

			var datatypesToDo []datatype.DataTypable
			// fetch all datatypes that have that property name
			for _, target := range datatypes {
				_, hasProperty := target.GetProperties()[propertyName]
				if hasProperty {
					datatypesToDo = append(datatypesToDo, target.GetProperties()[propertyName])
				}
			}

			UnifyUUID(datatypesToDo, idGenerator)

			propertiesDone[propertyName] = true
		}
	}
}

func SortScopes(input []*Scope) {
	sort.Slice(input, func(i, j int) bool {
		lineNumberA := input[i].Node.Source(false).StartLineNumber
		lineNumberB := input[j].Node.Source(false).StartLineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := input[i].Node.Source(false).StartColumnNumber
		columnNumberB := input[j].Node.Source(false).StartColumnNumber

		return *columnNumberA < *columnNumberB
	})
}
