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
	Type Type   `json:"type"`
	Name string `json:"name"`
}

type Variable struct {
	Name       string      `json:"string"`
	Complexity Complexity  `json:"complexity"`
	DataType   DataType    `json:"datatype"`
	Data       interface{} `json:"data"`
}
