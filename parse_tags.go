package encre

import (
	"fmt"
	"reflect"
	"strconv"
)

type mapping struct {
	NamedTags   map[string]string
	IndexedTags map[int]string
	Fields      map[string]string
}

func newMapping() *mapping {
	return &mapping{
		NamedTags:   make(map[string]string),
		IndexedTags: make(map[int]string),
		Fields:      make(map[string]string),
	}
}

func getMapping(o interface{}) (*mapping, error) {
	if o == nil {
		return nil, fmt.Errorf("object is nil")
	}

	v := reflect.ValueOf(o)

	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("object is not a pointer: %v", o)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("object is not a pointer to a struct: %v", o)
	}

	m := newMapping()

	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Type().Field(i)
		m.Fields[f.Name] = f.Name

		if t := f.Tag.Get("encre"); t != "" {
			n, err := strconv.Atoi(t)
			if err == nil {
				// indexed tag
				m.IndexedTags[n] = f.Name
			} else {
				// named tag
				m.NamedTags[t] = f.Name
			}
		}
	}

	return m, nil
}

var mappingCache = make(map[reflect.Type]*mapping)

func getCachedMapping(o interface{}) (*mapping, error) {
	t := reflect.TypeOf(o)
	m, cached := mappingCache[t]
	if cached {
		return m, nil
	}

	m, err := getMapping(o)
	if err != nil {
		return nil, err
	}

	mappingCache[t] = m

	return m, nil
}
