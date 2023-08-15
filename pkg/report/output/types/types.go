package types

import (
	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes    []dataflowtypes.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Risks        []dataflowtypes.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components   []dataflowtypes.Component    `json:"components,omitempty" yaml:"components,omitempty"`
	Dependencies []dataflowtypes.Dependency   `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Errors       []dataflowtypes.Error        `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type Output[T any] struct {
	Data         T
	Dataflow     *DataFlow
	Files        []string
	ReportFailed bool
}

type GenericOutput interface {
	ToGeneric() *Output[any]
}

func (output *Output[T]) ToGeneric() *Output[any] {
	if output == nil {
		return nil
	}

	return &Output[any]{
		Data:         output.Data,
		Dataflow:     output.Dataflow,
		Files:        output.Files,
		ReportFailed: output.ReportFailed,
	}
}

func ToSpecific[T any](generic *Output[any]) *Output[T] {
	if generic == nil {
		return nil
	}

	return &Output[T]{
		Data:         generic.Data.(T),
		Dataflow:     generic.Dataflow,
		Files:        generic.Files,
		ReportFailed: generic.ReportFailed,
	}
}
