package proto

//#cgo CXXFLAGS: -std=gnu++11
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_proto();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// Protobuf language using https://github.com/mitchellh/tree-sitter-proto
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_proto())
	return sitter.NewLanguage(ptr)
}
