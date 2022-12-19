package ffi

import (
	"fmt"
	"reflect"

	"github.com/open-policy-agent/opa/ast"
)

type opaqueTypeMap[T comparable] struct {
	typeName     string
	maxId        int
	idToInstance map[int]*T
	instanceToId map[*T]int
}

func newOpaqueTypeMap[T comparable](typeName string) opaqueTypeMap[T] {
	return opaqueTypeMap[T]{
		typeName:     "curio." + typeName,
		idToInstance: make(map[int]*T),
		instanceToId: make(map[*T]int),
	}
}

func (m *opaqueTypeMap[T]) CastToRegoInput(instance *T) interface{} {
	return map[string]interface{}{
		"type": m.typeName,
		"id":   m.getId(instance),
	}
}

func (m *opaqueTypeMap[T]) CastToRego(instance *T) *ast.Term {
	return ast.ObjectTerm(
		ast.Item(ast.StringTerm("type"), ast.StringTerm(m.typeName)),
		ast.Item(ast.StringTerm("id"), ast.IntNumberTerm(m.getId(instance))),
	)
}

func (m *opaqueTypeMap[T]) getId(instance *T) int {
	existingId, exists := m.instanceToId[instance]
	if exists {
		return existingId
	}

	id := m.maxId
	m.idToInstance[id] = instance
	m.instanceToId[instance] = id
	m.maxId = m.maxId + 1

	return id
}

func (m *opaqueTypeMap[T]) CastToGo(term *ast.Term) (*T, error) {
	object, ok := term.Value.(ast.Object)
	if !ok {
		return nil, fmt.Errorf("expected term to be object but was %s", reflect.TypeOf(term.Value))
	}

	typeName, err := getType(term, object)
	if err != nil {
		return nil, err
	}

	if typeName != m.typeName {
		return nil, fmt.Errorf("expected type '%s' but got '%s", m.typeName, typeName)
	}

	id, err := getId(term, object)
	if err != nil {
		return nil, err
	}

	node, ok := m.idToInstance[id]
	if !ok {
		return nil, fmt.Errorf("unknown instance with id %d", id)
	}

	return node, nil
}

func getType(term *ast.Term, object ast.Object) (string, error) {
	typeTerm := object.Get(ast.StringTerm("type"))
	if typeTerm == nil {
		return "", fmt.Errorf("expected 'type' attribute on term %s", term.String())
	}

	typeAst, ok := typeTerm.Value.(ast.String)
	if !ok {
		return "", fmt.Errorf("expected 'type' attribute to be a string on term %s", term.String())
	}

	return string(typeAst), nil
}

func getId(term *ast.Term, object ast.Object) (int, error) {
	idTerm := object.Get(ast.StringTerm("id"))
	if idTerm == nil {
		return 0, fmt.Errorf("expected 'id' attribute on term %s", term.String())
	}

	idAst, ok := idTerm.Value.(ast.Number)
	if !ok {
		return 0, fmt.Errorf("expected 'id' attribute to be a number on term %s", term.String())
	}

	id, ok := idAst.Int()
	if !ok {
		return 0, fmt.Errorf("expected 'id' attribute to be an int on term %s", term.String())
	}

	return id, nil
}
