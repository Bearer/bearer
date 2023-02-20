package types

import "github.com/bearer/bearer/new/detector/types"

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
