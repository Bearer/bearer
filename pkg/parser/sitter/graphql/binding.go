package graphql

//#cgo CXXFLAGS: -std=gnu++11
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_graphql();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// GraphQL language using https://github.com/bkegley/tree-sitter-graphql
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_graphql())
	return sitter.NewLanguage(ptr)
}
