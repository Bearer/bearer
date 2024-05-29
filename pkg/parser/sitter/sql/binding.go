package sql

//#cgo CFLAGS: -w
//#include <tree_sitter/parser.h>
//TSLanguage *tree_sitter_sql();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// SQL language using https://github.com/m-novikov/tree-sitter-sql
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_sql())
	return sitter.NewLanguage(ptr)
}
