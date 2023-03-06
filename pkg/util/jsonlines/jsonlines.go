package jsonlines

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/bearer/bearer/pkg/util/linescanner"
	"github.com/rs/zerolog/log"
)

const maxTokenSizeBytes int = 5 * 1024 * 1024

func getOriginalSlice(ptrToSlice interface{}) (slice reflect.Value, err error) {
	ptr2sl := reflect.TypeOf(ptrToSlice)
	if ptr2sl.Kind() != reflect.Ptr {
		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got %s", ptr2sl.Kind())
	}

	originalSlice := reflect.Indirect(reflect.ValueOf(ptrToSlice))
	sliceType := originalSlice.Type()
	if sliceType.Kind() != reflect.Slice {
		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got pointer to %s", sliceType.Kind())
	}
	return originalSlice, nil
}

func Encode(w io.Writer, ptrToSlice interface{}) error {
	slice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return fmt.Errorf("wrong value inputed into jsonlines encode: %s", err)
	}

	for i := 0; i < slice.Len(); i++ {
		err := json.NewEncoder(w).Encode(slice.Index(i).Interface())
		if err != nil {
			return fmt.Errorf("failed to encode json: %s", err)
		}
	}

	return nil
}

func Decode(r io.Reader, ptrToSlice interface{}) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return fmt.Errorf("wrong value inputed into jsonline decode: %s", err)
	}

	scanner := linescanner.NewSize(r, maxTokenSizeBytes)

	member := originalSlice.Type().Elem()

	for {
		ok := scanner.Scan()
		if !ok {
			break
		}

		item := scanner.Bytes()
		if len(item) == 0 {
			continue
		}

		newObj := reflect.New(member).Interface()

		log.Debug().Msgf("got bytes %s", string(item))

		err := json.Unmarshal(item, newObj)
		if err != nil {
			return fmt.Errorf("failed to unmarshal item: %s", err)
		}

		ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
		originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
	}

	return nil
}
