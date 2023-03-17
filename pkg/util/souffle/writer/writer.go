package writer

import "github.com/bearer/bearer/pkg/souffle/writer/base"

type FactWriter interface {
	Symbol(value string) base.Symbol
	Unsigned(value uint32) base.Unsigned
	Record(elements ...base.Element) base.Record
	WriteFact(relation string, elements ...base.Element) error
}
