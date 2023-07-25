package orm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func setupEncodeTests() (testApp *tests.TestApp, err error) {
	testApp, err = tests.NewTestApp()
	if err != nil {
		return nil, fmt.Errorf("could not create testApp: %w", err)
	}

	if err := testApp.Dao().SaveCollection(EntityWithAllPBTypes{}.Collection()); err != nil {
		return nil, fmt.Errorf("could not save collection: %w", err)
	}

	return
}

func TestEncode(t *testing.T) {
	testApp, err := setupEncodeTests()
	if err != nil {
		t.Fatalf("could not prepate esting environment: %v", err)
	}
	defer testApp.Cleanup()

	entity := entityExample
	actual, err := Encode(&entity, testApp.Dao())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := recordExample

	actualBusinessFields := RecordsColumnValueMap(actual)
	expectedBusinessFields := RecordsColumnValueMap(expected)

	if !reflect.DeepEqual(actualBusinessFields, expectedBusinessFields) {
		t.Errorf("expected %v, got %v", expectedBusinessFields, actualBusinessFields)
	}
}

func TestEncodeWithZeroValue(t *testing.T) {
	testApp, err := setupEncodeTests()
	if err != nil {
		t.Fatalf("could not prepate esting environment: %v", err)
	}
	defer testApp.Cleanup()

	entity := EntityWithAllPBTypes{}
	actual, err := Encode(&entity, testApp.Dao())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := models.NewRecord(EntityWithAllPBTypes{}.Collection())

	actualBusinessFields := RecordsColumnValueMap(actual)
	expectedBusinessFields := RecordsColumnValueMap(expected)

	if !reflect.DeepEqual(actualBusinessFields, expectedBusinessFields) {
		t.Errorf("expected %v, got %v", expectedBusinessFields, actualBusinessFields)
	}
}

func TestEncodeAll(t *testing.T) {
	testApp, err := setupEncodeTests()
	if err != nil {
		t.Fatalf("could not prepate esting environment: %v", err)
	}
	defer testApp.Cleanup()

	entity := EntityWithAllPBTypes{}

	actual, err := EncodeAll([]*EntityWithAllPBTypes{&entity, &entity}, testApp.Dao())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actualLen := len(actual); actualLen != 2 {
		t.Errorf("expected 2 elements, got %d", actualLen)
	}

	for _, actualElement := range actual {
		expected := models.NewRecord(EntityWithAllPBTypes{}.Collection())

		actualBusinessFields := RecordsColumnValueMap(actualElement)
		expectedBusinessFields := RecordsColumnValueMap(expected)

		if !reflect.DeepEqual(actualBusinessFields, expectedBusinessFields) {
			t.Errorf("expected %v, got %v", expectedBusinessFields, actualBusinessFields)
		}
	}
}
