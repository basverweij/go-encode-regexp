package encre

import (
	"reflect"
	"testing"
)

type test struct {
	Name       string
	NamedTag   string `encre:"namedTag"`
	IndexedTag string `encre:"1"`
}

func TestGetMapping(t *testing.T) {
	m, err := getMapping(nil)
	if m != nil || err == nil {
		t.Fatalf("getMapping(nil) returned non-nil mapping or nil error: mapping=%v, error=%v", m, err)
	}

	m, err = getMapping(test{})
	if m != nil || err == nil {
		t.Fatalf("getMapping(test{}) returned non-nil mapping or nil error: mapping=%v, error=%v", m, err)
	}

	m, err = getMapping(&test{})
	if m == nil || err != nil {
		t.Fatalf("getMapping(&test{}) returned nil mapping or non-nil error: mapping=%v, error=%v", m, err)
	}

	if m.Fields["Name"] != "Name" {
		t.Fatalf("wrong m.Fields mapping for 'Name': '%s' != 'Name'", m.Fields["Name"])
	}

	if m.NamedTags["namedTag"] != "NamedTag" {
		t.Fatalf("wrong m.NamedTags mapping for 'NamedTag': '%s' != 'NamedTag'", m.NamedTags["namedTag"])
	}

	if m.IndexedTags[1] != "IndexedTag" {
		t.Fatalf("wrong m.IndexedTags mapping for 'IndexedTag': '%s' != 'IndexedTag'", m.IndexedTags[1])
	}
}

func TestGetCachedMapping(t *testing.T) {
	m, err := getMapping(&test{})
	cm, cerr := getCachedMapping(&test{})

	if !reflect.DeepEqual(m, cm) {
		t.Fatalf("different result from cached mapping: %v != %v", m, cm)
	}

	if err != cerr {
		t.Fatalf("different error from cached mapping:\n%v != %v", err, cerr)
	}
}
