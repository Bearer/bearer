package types

import "github.com/bearer/bearer/new/detector/types"

type Object struct {
	Properties []Property
	// IsVirtual describes whether this object actually exists, or has
	// been detected as part of a variable name
	IsVirtual bool
}

type Property struct {
	Name   string
	Object *types.Detection
}

type String struct {
	Value     string
	IsLiteral bool
}
