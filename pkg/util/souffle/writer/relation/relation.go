package relationwriter

import (
	"github.com/bearer/bearer/pkg/souffle/binding"
	"github.com/bearer/bearer/pkg/util/souffle/writer/base"
)

type Writer struct {
	base.Base
	program *binding.Program
}

func New(program *binding.Program) *Writer {
	return &Writer{program: program}
}

func (writer *Writer) WriteFact(relationName string, elements ...base.Element) error {
	relation, err := writer.program.Relation(relationName)
	if err != nil {
		return err
	}

	tuple := relation.NewTuple()
	defer tuple.Close()

	for _, element := range elements {
		writer.writeElement(tuple, element)
	}

	relation.Insert(tuple)

	return nil
}

func (writer *Writer) writeElement(tuple *binding.Tuple, element base.Element) {
	switch value := element.(type) {
	case base.Symbol:
		tuple.WriteSymbol(string(value))
	case base.Unsigned:
		tuple.WriteUnsigned(uint32(value))
	case base.Record:
		tuple.WriteInteger(writer.packRecord(value))
	default:
		panic("unexpected element type")
	}
}

func (writer *Writer) packRecord(record base.Record) int32 {
	recordTuple := writer.program.NewRecordTuple(len(record))
	defer recordTuple.Close()

	for i, element := range record {
		switch value := element.(type) {
		case base.Symbol:
			recordTuple.WriteSymbol(i, string(value))
		case base.Unsigned:
			recordTuple.WriteUnsigned(i, uint32(value))
		case base.Record:
			recordTuple.WriteInteger(i, writer.packRecord(record))
		default:
			panic("unexpected element type")
		}
	}

	return writer.program.PackRecord(recordTuple)
}
