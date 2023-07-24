package orm

import (
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestDecode(t *testing.T) {
	r := recordExample

	entity := EntityWithAllPBTypes{}
	if err := Decode(r, &entity); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := entityExample
	if !reflect.DeepEqual(entity, expected) {
		t.Errorf("expected %v, got %v", expected, entity)
	}
}

func TestDecodeAll(t *testing.T) {
	records := []*models.Record{recordExample, recordExample}
	entities := []*EntityWithAllPBTypes{{}, {}}

	if err := DecodeAll(records, entities); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	for _, entity := range entities {
		if entity == nil {
			t.Errorf("expected non-nil pointer, got nil")
			continue
		}

		expected := entityExample
		if !reflect.DeepEqual(*entity, expected) {
			t.Errorf("expected %v, got %v", expected, *entity)
		}
	}

}
