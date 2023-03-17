#include "binding.h"
#include "souffle/SouffleInterface.h"

extern "C"
{
  SouffleProgram souffle_program_factory_new_instance(const char *name)
  {
    return souffle::ProgramFactory::newInstance(name);
  }
}

typedef struct
{
  SouffleSymbolTable symbolTable;
  int32_t *elements;
  size_t arity;
  bool owned;
} InternalRecordTuple;

extern "C"
{
  SouffleRecordTuple souffle_program_new_record_tuple(const SouffleProgram program, size_t arity)
  {
    InternalRecordTuple *tuple = new InternalRecordTuple;

    tuple->symbolTable = &((souffle::SouffleProgram *)program)->getSymbolTable();
    tuple->arity = arity;
    tuple->elements = (int32_t *)malloc(arity * sizeof(int32_t));
    tuple->owned = true;

    return tuple;
  }

  SouffleRelation souffle_program_get_relation(const SouffleProgram program, const char *name)
  {
    return ((souffle::SouffleProgram *)program)->getRelation(name);
  }

  int32_t souffle_program_pack_record(const SouffleProgram program, const SouffleRecordTuple tuple)
  {
    souffle::RecordTable &recordTable = ((souffle::SouffleProgram *)program)->getRecordTable();
    InternalRecordTuple *internalTuple = (InternalRecordTuple *)tuple;

    return recordTable.pack(internalTuple->elements, internalTuple->arity);
  }

  SouffleRecordTuple souffle_program_unpack_record(const SouffleProgram program, int32_t index, size_t arity)
  {
    souffle::SouffleProgram *souffleProgram = (souffle::SouffleProgram *)program;

    InternalRecordTuple *tuple = new InternalRecordTuple;
    tuple->symbolTable = &souffleProgram->getSymbolTable();
    tuple->arity = arity;
    tuple->elements = (int32_t *)souffleProgram->getRecordTable().unpack(index, arity);
    tuple->owned = false;

    return tuple;
  }

  void souffle_program_run(const SouffleProgram program)
  {
    ((souffle::SouffleProgram *)program)->run();
  }

  void souffle_program_free(SouffleProgram program)
  {
    delete (souffle::SouffleProgram *)program;
  }
}

extern "C"
{
  const char *souffle_record_tuple_read_symbol(const SouffleRecordTuple tuple, size_t index)
  {
    InternalRecordTuple *internalTuple = (InternalRecordTuple *)tuple;
    souffle::SymbolTable &symbolTable = (souffle::SymbolTable &)internalTuple->symbolTable;

    return symbolTable.decode(index).c_str();
  }

  uint32_t souffle_record_tuple_read_unsigned(const SouffleRecordTuple tuple, size_t index)
  {
    InternalRecordTuple *internalTuple = (InternalRecordTuple *)tuple;

    return souffle::ramBitCast<souffle::RamUnsigned>(internalTuple->elements[index]);
  }

  int32_t souffle_record_tuple_read_integer(const SouffleRecordTuple tuple, size_t index)
  {
    return ((InternalRecordTuple *)tuple)->elements[index];
  }

  void souffle_record_tuple_write_symbol(const SouffleRecordTuple tuple, size_t index, const char *value)
  {
    InternalRecordTuple *internalTuple = (InternalRecordTuple *)tuple;
    souffle::SymbolTable &symbolTable = (souffle::SymbolTable &)internalTuple->symbolTable;

    internalTuple->elements[index] = symbolTable.encode(value);
  }

  void souffle_record_tuple_write_unsigned(const SouffleRecordTuple tuple, size_t index, uint32_t value)
  {
    ((InternalRecordTuple *)tuple)->elements[index] = souffle::ramBitCast(value);
  }

  void souffle_record_tuple_write_integer(const SouffleRecordTuple tuple, size_t index, int32_t value)
  {
    ((InternalRecordTuple *)tuple)->elements[index] = value;
  }

  void souffle_record_tuple_free(SouffleRecordTuple tuple)
  {
    InternalRecordTuple *internalTuple = (InternalRecordTuple *)tuple;

    if (internalTuple->owned) {
      free(internalTuple->elements);
    }

    free((void *)tuple);
  }
}

typedef struct
{
  souffle::Relation::iterator current;
  souffle::Relation::iterator end;
} InternalRelationIterator;

extern "C"
{
  SouffleTuple souffle_relation_new_tuple(const SouffleRelation relation)
  {
    return new souffle::tuple((souffle::Relation *)relation);
  }

  SouffleRelationIterator souffle_relation_new_iterator(const SouffleRelation relation)
  {
    souffle::Relation *souffleRelation = (souffle::Relation *)relation;

    return new InternalRelationIterator{souffleRelation->begin(), souffleRelation->end()};
  }

  void souffle_relation_insert(const SouffleRelation relation, const SouffleTuple tuple)
  {
    ((souffle::Relation *)relation)->insert(*(souffle::tuple *)tuple);
  }

  size_t souffle_relation_size(const SouffleRelation relation)
  {
    return ((souffle::Relation *)relation)->size();
  }

  size_t souffle_relation_arity(const SouffleRelation relation)
  {
    return ((souffle::Relation *)relation)->getArity();
  }

  const char *souffle_relation_attr_type(const SouffleRelation relation, size_t index)
  {
    return ((souffle::Relation *)relation)->getAttrType(index);
  }
}

extern "C"
{
  char *souffle_tuple_read_symbol(const SouffleTuple tuple)
  {
    std::string value;
    (*(souffle::tuple *)tuple) >> value;

    const size_t bufferSize = value.size() + 1; // including null byte
    char *buffer = (char *)malloc(bufferSize * sizeof(char));
    memcpy(buffer, value.c_str(), bufferSize);

    return buffer;
  }

  uint32_t souffle_tuple_read_unsigned(const SouffleTuple tuple)
  {
    uint32_t value;
    (*(souffle::tuple *)tuple) >> value;

    return value;
  }

  int32_t souffle_tuple_read_integer(const SouffleTuple tuple)
  {
    int32_t value;
    (*(souffle::tuple *)tuple) >> value;

    return value;
  }

  void souffle_tuple_write_symbol(const SouffleTuple tuple, const char *value)
  {
    std::string stringValue(value);
    (*(souffle::tuple *)tuple) << stringValue;
  }

  void souffle_tuple_write_unsigned(const SouffleTuple tuple, uint32_t value)
  {
    (*(souffle::tuple *)tuple) << value;
  }

  void souffle_tuple_write_integer(const SouffleTuple tuple, int32_t value)
  {
    (*(souffle::tuple *)tuple) << value;
  }

  void souffle_tuple_free(SouffleTuple tuple)
  {
    delete (souffle::tuple *)tuple;
  }
}

extern "C"
{
  bool souffle_relation_iterator_has_next(const SouffleRelationIterator iterator)
  {
    InternalRelationIterator *internalIterator = (InternalRelationIterator *)iterator;

    return internalIterator->current != internalIterator->end;
  }

  SouffleTuple souffle_relation_iterator_get_next(const SouffleRelationIterator iterator)
  {
    InternalRelationIterator *internalIterator = (InternalRelationIterator *)iterator;

    SouffleTuple result = &*internalIterator->current;
    internalIterator->current++;

    return result;
  }

  void souffle_relation_iterator_free(SouffleRelationIterator iterator)
  {
    InternalRelationIterator *internalIterator = (InternalRelationIterator *)iterator;

    delete internalIterator;
  }
}

extern "C"
{
  void c_free(void *c)
  {
    free(c);
  }
}
