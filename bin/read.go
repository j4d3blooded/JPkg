package jpkg_bin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func BinaryRead[T any](r io.Reader) (*T, error) {
	rt := reflect.TypeFor[T]()
	rv := reflect.New(rt).Elem()

	if rt.Kind() != reflect.Struct {
		return nil, errors.New("Binary Write only works with structs")
	}

	for fI := range rv.NumField() {
		if !rt.Field(fI).IsExported() {
			continue
		}

		field := rv.Field(fI)

		if field.Kind() == reflect.Struct {
			return nil, errors.New("struct fields are not supported")
		}

		if field.Kind() == reflect.String { // write string as sized utf8

			var length uint64

			if err := br(r, &length); err != nil {
				return nil, fmt.Errorf("error reading string field %v length: %w", fI, err)
			}

			str := make([]rune, length)

			if err := br(r, &str); err != nil {
				return nil, fmt.Errorf("error reading string field %v: %w", fI, err)
			}

			field.SetString(string(str))
			continue
		}

		fv := reflect.New(field.Type()).Interface()

		if err := br(r, fv); err != nil {
			return nil, fmt.Errorf("error reading field %v: %w", fI, err)
		}

		field.Set(reflect.ValueOf(fv).Elem())
	}

	a := rv.Interface().(T)

	return &a, nil
}

func br(r io.Reader, val any) error {
	return binary.Read(r, binary.BigEndian, val)
}
