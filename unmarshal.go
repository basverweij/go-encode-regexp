package encre

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
)

type Unmarshaler interface {
	UnmarshalRegExp(s string) error
}

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func Unmarshal(data []byte, v interface{}, expr string) error {
	d, err := NewDecoder(bytes.NewBuffer(data), expr)
	if err != nil {
		return err
	}

	return d.Decode(v)
}

type decoder struct {
	r  io.RuneReader
	re *regexp.Regexp
}

func NewDecoder(r io.Reader, expr string) (*decoder, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &decoder{r: bufio.NewReader(r), re: re}, nil
}

func (d *decoder) Decode(v interface{}) error {
	m, err := getCachedMapping(v)
	if err != nil {
		return err
	}

	c := NewCloningRuneReader(d.r)
	sm := d.re.FindReaderSubmatchIndex(c)
	if sm == nil || len(sm) == 0 {
		return fmt.Errorf("expression not found in input")
	}

	smCount := d.re.NumSubexp()
	smNames := d.re.SubexpNames()

	// loop all sub matches
	o := reflect.ValueOf(v).Elem()
	for i := 1; i <= smCount; i++ {
		if sm[i*2] == -1 {
			// subgroup not matched
			continue
		}

		name := ""
		if n := m.NamedTags[smNames[i]]; n != "" {
			name = n // use mapped named tag
		} else if n := m.IndexedTags[i]; n != "" {
			name = n // use mapped indexed tag
		} else if n := m.Fields[smNames[i]]; n != "" {
			name = n // use field name
		}

		if name == "" {
			// sub match not mapped to struct field
			continue
		}

		f := o.FieldByName(name)
		err := setField(&f, c.Slice(sm[i*2], sm[i*2+1]))
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	notSupported int = iota
	boolType
	intType
	uintType
	floatType
	stringType
)

func setField(f *reflect.Value, s string) error {
	if f.CanInterface() && f.Type().Implements(unmarshalerType) {
		return f.Interface().(Unmarshaler).UnmarshalRegExp(s)
	}

	if f.CanInterface() && f.Type().Implements(textUnmarshalerType) {
		return f.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
	}

	t := notSupported
	size := 0

	switch f.Kind() {
	case reflect.Bool:
		t = boolType
	case reflect.Int:
		t, size = intType, 64
	case reflect.Int8:
		t, size = intType, 8
	case reflect.Int16:
		t, size = intType, 16
	case reflect.Int32:
		t, size = intType, 32
	case reflect.Int64:
		t, size = intType, 64
	case reflect.Uint:
		t, size = uintType, 64
	case reflect.Uint8:
		t, size = uintType, 8
	case reflect.Uint16:
		t, size = uintType, 16
	case reflect.Uint32:
		t, size = uintType, 32
	case reflect.Uint64:
		t, size = uintType, 64
	case reflect.Float32:
		t, size = floatType, 32
	case reflect.Float64:
		t, size = floatType, 64
	case reflect.String:
		t = stringType
	}

	if t == notSupported {
		return fmt.Errorf("unsupported type: %v", f.Type())
	}

	switch t {
	case boolType:
		x, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		f.SetBool(x)
	case intType:
		x, err := strconv.ParseInt(s, 10, size)
		if err != nil {
			return err
		}
		f.SetInt(x)
	case uintType:
		x, err := strconv.ParseUint(s, 10, size)
		if err != nil {
			return err
		}
		f.SetUint(x)
	case floatType:
		x, err := strconv.ParseFloat(s, size)
		if err != nil {
			return err
		}
		f.SetFloat(x)
	case stringType:
		f.SetString(s)
	}

	return nil
}
