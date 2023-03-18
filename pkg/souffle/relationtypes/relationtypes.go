package relationtypes

/* This package defines Go types that map to Souffle relations and records.
 * The order of fields must match the order in the Souffle definition.
 */

type Rule_Match struct {
	RuleName     string
	PatternIndex uint32
	Node         uint32
	Location     AST_Location
}

type AST_Location struct {
	StartByte   uint32
	StartLine   uint32
	StartColumn uint32
	EndLine     uint32
	EndColumn   uint32
}
