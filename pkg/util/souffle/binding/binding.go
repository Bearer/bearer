package binding

// #cgo CXXFLAGS: -std=c++17
// #include "binding.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Program struct {
	c C.SouffleProgram
}

type Relation struct {
	c C.SouffleRelation
}

type RelationIterator struct {
	c        C.SouffleRelationIterator
	relation *Relation
}

type Tuple struct {
	c        C.SouffleTuple
	relation *Relation
	owned    bool
}

type RecordTuple struct {
	c C.SouffleRecordTuple
}

// A large number but low enough to be supported on all arch's
const maxCArrayIndex = 1<<30 - 1

func NewProgram(name string) (*Program, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cProgram := C.souffle_program_factory_new_instance(cName)
	if cProgram == nil {
		return nil, fmt.Errorf("%s program not found", name)
	}

	program := &Program{c: cProgram}
	runtime.SetFinalizer(program, func(finalizedProgram *Program) {
		finalizedProgram.Close()
	})

	return program, nil
}

func (program *Program) Relation(name string) (*Relation, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cRelation := C.souffle_program_get_relation(program.c, cName)
	if cRelation == nil {
		return nil, fmt.Errorf("%s relation not found", name)
	}

	return &Relation{c: cRelation}, nil
}

func (program *Program) PackRecord(tuple *RecordTuple) int32 {
	return int32(C.souffle_program_pack_record(program.c, tuple.c))
}

func (program *Program) UnpackRecord(index int32, arity int) (tuple *RecordTuple) {
	cTuple := C.souffle_program_unpack_record(program.c, C.int32_t(index), C.size_t(arity))

	return &RecordTuple{c: cTuple}
}

func (program *Program) NewRecordTuple(arity int) *RecordTuple {
	cTuple := C.souffle_program_new_record_tuple(program.c, C.size_t(arity))

	tuple := &RecordTuple{c: cTuple}
	runtime.SetFinalizer(tuple, func(finalizedTuple *RecordTuple) {
		finalizedTuple.Close()
	})

	return tuple
}

func (program *Program) Run() {
	C.souffle_program_run(program.c)
}

func (program *Program) Close() {
	runtime.SetFinalizer(program, nil)
	C.souffle_program_free(program.c)
}

func (relation *Relation) NewIterator() *RelationIterator {
	cIterator := C.souffle_relation_new_iterator(relation.c)

	iterator := &RelationIterator{c: cIterator, relation: relation}
	runtime.SetFinalizer(iterator, func(finalizedIterator *RelationIterator) {
		finalizedIterator.Close()
	})

	return iterator
}

func (relation *Relation) NewTuple() *Tuple {
	cTuple := C.souffle_relation_new_tuple(relation.c)

	tuple := &Tuple{c: cTuple, relation: relation, owned: true}
	runtime.SetFinalizer(tuple, func(finalizedTuple *Tuple) {
		finalizedTuple.Close()
	})

	return tuple
}

func (relation *Relation) Insert(tuple *Tuple) {
	C.souffle_relation_insert(relation.c, tuple.c)
}

func (relation *Relation) Size() int {
	return int(C.souffle_relation_size(relation.c))
}

func (relation *Relation) Arity() int {
	return int(C.souffle_relation_arity(relation.c))
}

func (relation *Relation) AttrType(index int) string {
	cValue := C.souffle_relation_attr_type(relation.c, C.size_t(index))
	return C.GoString(cValue)
}

func (iterator *RelationIterator) HasNext() bool {
	return bool(C.souffle_relation_iterator_has_next(iterator.c))
}

func (iterator *RelationIterator) GetNext() *Tuple {
	cTuple := C.souffle_relation_iterator_get_next(iterator.c)

	return &Tuple{c: cTuple, relation: iterator.relation, owned: false}
}

func (iterator *RelationIterator) Close() {
	runtime.SetFinalizer(iterator, nil)
	C.souffle_relation_iterator_free(iterator.c)
}

func (tuple *Tuple) ReadSymbol() string {
	cValue := C.souffle_tuple_read_symbol(tuple.c)
	defer C.free(unsafe.Pointer(cValue))

	return C.GoString(cValue)
}

func (tuple *Tuple) ReadUnsigned() uint32 {
	return uint32(C.souffle_tuple_read_unsigned(tuple.c))
}

func (tuple *Tuple) ReadInteger() int32 {
	return int32(C.souffle_tuple_read_integer(tuple.c))
}

func (tuple *Tuple) WriteSymbol(value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	C.souffle_tuple_write_symbol(tuple.c, cValue)
}

func (tuple *Tuple) WriteUnsigned(value uint32) {
	C.souffle_tuple_write_unsigned(tuple.c, C.uint32_t(value))
}

func (tuple *Tuple) WriteInteger(value int32) {
	C.souffle_tuple_write_integer(tuple.c, C.int32_t(value))
}

func (tuple *Tuple) Relation() *Relation {
	return tuple.relation
}

func (tuple *Tuple) Close() {
	if tuple.owned {
		runtime.SetFinalizer(tuple, nil)
		C.souffle_tuple_free(tuple.c)
	}
}

func (tuple *RecordTuple) ReadSymbol(index int) string {
	return C.GoString(C.souffle_record_tuple_read_symbol(tuple.c, C.size_t(index)))
}

func (tuple *RecordTuple) ReadUnsigned(index int) uint32 {
	return uint32(C.souffle_record_tuple_read_unsigned(tuple.c, C.size_t(index)))
}

func (tuple *RecordTuple) ReadInteger(index int) int32 {
	return int32(C.souffle_record_tuple_read_integer(tuple.c, C.size_t(index)))
}

func (tuple *RecordTuple) WriteSymbol(index int, value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	C.souffle_record_tuple_write_symbol(tuple.c, C.size_t(index), cValue)
}

func (tuple *RecordTuple) WriteUnsigned(index int, value uint32) {
	C.souffle_record_tuple_write_unsigned(tuple.c, C.size_t(index), C.uint32_t(value))
}

func (tuple *RecordTuple) WriteInteger(index int, value int32) {
	C.souffle_record_tuple_write_integer(tuple.c, C.size_t(index), C.int32_t(value))
}

func (tuple *RecordTuple) Close() {
	runtime.SetFinalizer(tuple, nil)
	C.souffle_record_tuple_free(tuple.c)
}
