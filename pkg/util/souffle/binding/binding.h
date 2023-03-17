#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef void *SouffleProgram;
typedef void *SouffleRelation;
typedef void *SouffleTuple;
typedef void *SouffleSymbolTable;
typedef void *SouffleRelationIterator;
typedef void *SouffleRecordTuple;

#ifdef __cplusplus
extern "C"
{
#endif
  SouffleProgram souffle_program_factory_new_instance(const char *name);

  SouffleRecordTuple souffle_program_new_record_tuple(const SouffleProgram program, size_t arity);
  SouffleRelation souffle_program_get_relation(const SouffleProgram program, const char *name);
  int32_t souffle_program_pack_record(const SouffleProgram program, const SouffleRecordTuple tuple);
  SouffleRecordTuple souffle_program_unpack_record(const SouffleProgram program, int32_t index, size_t arity);
  void souffle_program_run(const SouffleProgram program);
  void souffle_program_free(SouffleProgram program);

  SouffleTuple souffle_relation_new_tuple(const SouffleRelation relation);
  SouffleRelationIterator souffle_relation_new_iterator(const SouffleRelation relation);
  void souffle_relation_insert(const SouffleRelation relation, const SouffleTuple tuple);
  size_t souffle_relation_size(const SouffleRelation relation);
  size_t souffle_relation_arity(const SouffleRelation relation);
  const char *souffle_relation_attr_type(const SouffleRelation relation, size_t index);

  char *souffle_tuple_read_symbol(const SouffleTuple tuple);
  uint32_t souffle_tuple_read_unsigned(const SouffleTuple tuple);
  int32_t souffle_tuple_read_integer(const SouffleTuple tuple);
  void souffle_tuple_write_symbol(const SouffleTuple tuple, const char *value);
  void souffle_tuple_write_unsigned(const SouffleTuple tuple, uint32_t value);
  void souffle_tuple_write_integer(const SouffleTuple tuple, int32_t value);
  void souffle_tuple_free(SouffleTuple tuple);

  const char *souffle_record_tuple_read_symbol(const SouffleRecordTuple tuple, size_t index);
  uint32_t souffle_record_tuple_read_unsigned(const SouffleRecordTuple tuple, size_t index);
  int32_t souffle_record_tuple_read_integer(const SouffleRecordTuple tuple, size_t index);
  void souffle_record_tuple_write_symbol(const SouffleRecordTuple tuple, size_t index, const char *value);
  void souffle_record_tuple_write_unsigned(const SouffleRecordTuple tuple, size_t index, uint32_t value);
  void souffle_record_tuple_write_integer(const SouffleRecordTuple tuple, size_t index, int32_t value);
  void souffle_record_tuple_free(SouffleRecordTuple tuple);

  bool souffle_relation_iterator_has_next(const SouffleRelationIterator iterator);
  SouffleTuple souffle_relation_iterator_get_next(const SouffleRelationIterator iterator);
  void souffle_relation_iterator_free(SouffleRelationIterator iterator);
#ifdef __cplusplus
}
#endif
