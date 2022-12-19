package ffi

import (
	"reflect"

	"github.com/bearer/curio/new/language"
	"github.com/open-policy-agent/opa/ast"
)

func interfaceToValue(x interface{}) (ast.Value, error) {
	typ := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if node, isNode := x.(*language.Node); isNode {
		data, err := GetData()
		if err != nil {
			return nil, err
		}

		return data.nodes.CastToRego(node).Value, nil
	}

	switch typ.Kind() {
	case reflect.Struct:
		items := make([][2]*ast.Term, typ.NumField())

		for i := 0; i < typ.NumField(); i++ {
			fieldValue, err := interfaceToValue(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}

			items[i] = ast.Item(
				ast.StringTerm(typ.Field(i).Name),
				ast.NewTerm(fieldValue),
			)
		}

		return ast.NewObject(items...), nil
	case reflect.Array, reflect.Slice:
		elements := make([]*ast.Term, v.Len())

		for i := 0; i < v.Len(); i++ {
			elementValue, err := interfaceToValue(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}

			elements[i] = ast.NewTerm(elementValue)
		}

		// json.Marshal()

		return ast.NewArray(elements...), nil
	default:
		return ast.InterfaceToValue(x)
	}
}
