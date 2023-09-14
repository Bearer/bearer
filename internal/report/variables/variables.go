package variables

type Type string
type Complexity string
type DataType string

const (
	VariableComplexitySimple Complexity = "simple"
	VariableComplexityObject Complexity = "object"
)
const (
	VariableDataTypeString  DataType = "string"
	VariableDataTypeNumber  DataType = "number"
	VariableDataTypeBoolean DataType = "boolean"
)

const (
	VariableEnvironment Type = "environment"
	VariableTemplate    Type = "template"
	VariableName        Type = "variable"
)

type Identifier struct {
	Type Type   `json:"type" yaml:"type"`
	Name string `json:"name" yaml:"name"`
}

type Variable struct {
	Name       string      `json:"string" yaml:"string"`
	Complexity Complexity  `json:"complexity" yaml:"complexity"`
	DataType   DataType    `json:"datatype" yaml:"datatype"`
	Data       interface{} `json:"data" yaml:"data"`
}
