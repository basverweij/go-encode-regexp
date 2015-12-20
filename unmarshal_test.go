package encre

import (
	"reflect"
	"testing"
)

var format = `(?i)^IndexedTag: ([0-9]+), Name: (?P<Name>[a-z]+), NamedTag: (?P<namedTag>[a-z]+)$`
var data = []byte("IndexedTag: 123, Name: abc, NamedTag: xyz")

func TestUnmarshal(t *testing.T) {
	err := Unmarshal(data, nil, format)
	if err == nil {
		t.Fatalf("Unmarshal(nil) should return error")
	}

	err = Unmarshal(data, test{}, format)
	if err == nil {
		t.Fatalf("Unmarshal(test{}) should return error")
	}

	actual := &test{}
	err = Unmarshal(data, actual, format)
	if err != nil {
		t.Fatalf("Unmarshal(&test{}) should not return error")
	}

	expected := &test{IndexedTag: 123, Name: "abc", NamedTag: "xyz"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unmarshalled actual result not equal to expected: %v != %v", actual, expected)
	}
}
