package jpkg_bin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func BinaryWrite(w io.Writer, data any) error {

	rv := reflect.ValueOf(data)

	if rv.Kind() != reflect.Struct {
		return errors.New("Binary Write only works with structs")
	}

	rt := rv.Type()

	for fI := range reflect.ValueOf(data).NumField() {
		if !rt.Field(fI).IsExported() {
			continue
		}

		field := rv.Field(fI)
		value := field.Interface()

		if field.Kind() == reflect.Struct { // handle embeded structs
			return errors.New("struct fields are not supported")
		}

		if field.Kind() == reflect.String { // write string as sized utf8

			value := value.(string)

			if err := bw(w, uint64(len(value))); err != nil {
				return fmt.Errorf("error writing string field %v length: %w", fI, err)
			}

			if err := bw(w, []rune(value)); err != nil {
				return fmt.Errorf("error writing string field %v: %w", fI, err)
			}

			continue
		}

		if err := bw(w, value); err != nil {
			return fmt.Errorf("error writing field %v: %w", fI, err)
		}
	}

	return nil
}

func bw(w io.Writer, data any) error {
	return binary.Write(w, binary.BigEndian, data)
}
