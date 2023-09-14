package toml

//#cgo CXXFLAGS: -std=gnu++11
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_toml_1();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// Protobuf language using https://github.com/ikatyang/tree-sitter-toml (check readme)
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_toml_1())
	return sitter.NewLanguage(ptr)
}
