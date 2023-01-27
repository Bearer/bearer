package types

import "github.com/bearer/curio/new/detector/types"

type Object struct {
	Name       string
	Properties []*types.Detection
}

type Property struct {
	Name string
}

type String struct {
	Value string
}
