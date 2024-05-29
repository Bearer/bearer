package config_variables

//#cgo CXXFLAGS: -std=gnu++11
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_config_variables();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_config_variables())
	return sitter.NewLanguage(ptr)
}
