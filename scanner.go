package sqlscan

import (
	"reflect"
)

type Scanner interface {
	Scan(v any) error
}

type Scannable interface {
	Scan(dest ...any) error
	Columns() ([]string, error)
}

type StructScanner struct {
	src Scannable
}

func New(src Scannable) Scanner {
	return StructScanner{
		src: src,
	}
}

func (s StructScanner) Scan(v any) error {
	columns, err := s.src.Columns()
	if err != nil {
		return err
	}
	fields := make([]any, 0)
	for _, column := range columns {
		field, ok := s.getFieldValue(column, v)
		if !ok {
			// There is no field with tag matching column name, ignore and proceed to next column
			continue
		}
		fields = append(fields, field)
	}
	return s.src.Scan(fields...)
}

func (s StructScanner) getFieldValue(tag string, v any) (any, bool) {
	ref := reflect.Indirect(reflect.ValueOf(v))
	for curr := 0; curr < ref.NumField(); curr++ {
		addr := ref.Field(curr).Addr()
		field := ref.Type().Field(curr)
		if field.Tag.Get("sql") == tag {
			return addr.Interface(), true
		}
	}
	return nil, false
}
