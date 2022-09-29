package xml2

//#cgo CXXFLAGS: -std=gnu++11
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_xml2();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// tree sitter using https://github.com/dorgnarg/tree-sitter-xml (check readme)
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_xml2())
	return sitter.NewLanguage(ptr)
}
