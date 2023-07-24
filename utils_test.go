package orm

import (
	"reflect"
	"testing"
)

func TestPointer(t *testing.T) {
	dataset := []struct {
		label string
		input interface{}
	}{
		{
			label: "with string",
			input: "foo",
		},
		{
			label: "with int",
			input: 12,
		},
		{
			label: "with map",
			input: map[string]int{"foo": 12},
		},
	}

	for _, tt := range dataset {
		t.Run(tt.label, func(t *testing.T) {
			expected := &tt.input
			actual := pointer(tt.input)
			if !reflect.DeepEqual(expected, actual) {
				t.Error("unexpected pointed value")
			}
		})
	}
}
