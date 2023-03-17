package souffle

import (
	"fmt"
	"reflect"

	_ "github.com/bearer/bearer/pkg/souffle/rules" // FIXME: this should be up in main or something
	"github.com/bearer/bearer/pkg/util/souffle/binding"
)

type Souffle struct {
	program *binding.Program
}

func New(programName string) (*Souffle, error) {
	program, err := binding.NewProgram(programName)
	if err != nil {
		return nil, err
	}

	return &Souffle{program: program}, nil
}

func (souffle *Souffle) Program() *binding.Program {
	return souffle.program
}

func (souffle *Souffle) Run() {
	souffle.program.Run()
}

func (souffle *Souffle) Relation(name string) (*binding.Relation, error) {
	return souffle.program.Relation(name)
}

func (souffle *Souffle) Marshal(relation *binding.Relation, value any) (*binding.Tuple, error) {
	fields, err := getFields(reflect.TypeOf(value))
	if err != nil {
		return nil, err
	}

	if relation.Arity() != len(fields) {
		return nil, fmt.Errorf("mismatch in arity (relation=%d, struct=%d)", relation.Arity(), len(fields))
	}

	typeValue := reflect.ValueOf(value)
	tuple := relation.NewTuple()

	for _, field := range fields {
		fieldValue := typeValue.FieldByIndex(field.Index)

		switch field.Type.Kind() {
		case reflect.String:
			tuple.WriteSymbol(fieldValue.String())
		case reflect.Uint32:
			tuple.WriteUnsigned(uint32(fieldValue.Uint()))
		case reflect.Int32:
			tuple.WriteInteger(int32(fieldValue.Int()))
		case reflect.Struct:
			index, err := souffle.encodeRecord(fieldValue)
			if err != nil {
				return nil, fmt.Errorf("encoding error for field %s: %w", field.Name, err)
			}

			tuple.WriteInteger(index)
		default:
			return nil, fmt.Errorf("kind %s not supported for field %s", field.Type.Kind(), field.Name)
		}
	}

	return tuple, nil
}

func (souffle *Souffle) Unmarshal(destination any, tuple *binding.Tuple) error {
	typ := reflect.TypeOf(destination)
	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("type %s not a pointer but %s", typ.Name(), typ.Kind())
	}

	fields, err := getFields(typ.Elem())
	if err != nil {
		return err
	}

	if tuple.Relation().Arity() != len(fields) {
		return fmt.Errorf("mismatch in arity (tuple=%d, struct=%d)", tuple.Relation().Arity(), len(fields))
	}

	typeValue := reflect.Indirect(reflect.ValueOf(destination))

	for _, field := range fields {

		fieldValue := typeValue.FieldByIndex(field.Index)

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(tuple.ReadSymbol())
		case reflect.Uint32:
			fieldValue.SetUint(uint64(tuple.ReadUnsigned()))
		case reflect.Int32:
			fieldValue.SetInt(int64(tuple.ReadInteger()))
		case reflect.Struct:
			err := souffle.decodeRecord(field.Type, fieldValue.Addr(), tuple.ReadInteger())
			if err != nil {
				return fmt.Errorf("decoding error for field %s: %w", field.Name, err)
			}
		default:
			return fmt.Errorf("kind %s not supported (field %s)", field.Type.Kind(), field.Name)
		}
	}

	return nil
}

func (souffle *Souffle) Close() {
	souffle.program.Close()
}

func (souffle *Souffle) encodeRecord(value reflect.Value) (int32, error) {
	fields, err := getFields(value.Type())
	if err != nil {
		return 0, err
	}

	// TODO: Is there some way to check the arity?

	recordTuple := souffle.program.NewRecordTuple(len(fields))
	defer recordTuple.Close()

	for i, field := range fields {
		fieldValue := value.FieldByIndex(field.Index)

		switch field.Type.Kind() {
		case reflect.String:
			recordTuple.WriteSymbol(i, fieldValue.String())
		case reflect.Uint32:
			recordTuple.WriteUnsigned(i, uint32(fieldValue.Uint()))
		case reflect.Int32:
			recordTuple.WriteInteger(i, int32(fieldValue.Int()))
		case reflect.Struct:
			index, err := souffle.encodeRecord(fieldValue)
			if err != nil {
				return 0, fmt.Errorf("encoding error for field %s: %w", field.Name, err)
			}

			recordTuple.WriteInteger(i, index)
		default:
			return 0, fmt.Errorf("kind %s not supported for field %s", field.Type.Kind(), field.Name)
		}
	}

	return souffle.program.PackRecord(recordTuple), nil
}

func (souffle *Souffle) decodeRecord(typ reflect.Type, destination reflect.Value, index int32) error {
	fields, err := getFields(typ)
	if err != nil {
		return err
	}

	// TODO: Is there some way to check the arity?

	recordTuple := souffle.program.UnpackRecord(index, len(fields))
	defer recordTuple.Close()

	value := reflect.Indirect(destination)

	for i, field := range fields {
		fieldValue := value.FieldByIndex(field.Index)

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(recordTuple.ReadSymbol(i))
		case reflect.Uint32:
			fieldValue.SetUint(uint64(recordTuple.ReadUnsigned(i)))
		case reflect.Int32:
			fieldValue.SetInt(int64(recordTuple.ReadInteger(i)))
		case reflect.Struct:
			err := souffle.decodeRecord(field.Type, fieldValue.Addr(), recordTuple.ReadInteger(i))
			if err != nil {
				return fmt.Errorf("decoding error for field %s: %w", field.Name, err)
			}
		default:
			return fmt.Errorf("kind %s not supported for field %s", field.Type.Kind(), field.Name)
		}
	}

	return nil
}

func getFields(typ reflect.Type) ([]reflect.StructField, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s not a struct but %s", typ.Name(), typ.Kind())
	}

	var fields []reflect.StructField

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("souffle")

		if field.IsExported() && tag != "-" {
			fields = append(fields, field)
		}
	}

	return fields, nil
}
