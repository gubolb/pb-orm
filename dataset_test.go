package orm

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ Entity = EntityWithAllPBTypes{}

type StringUnderlyingType string

const (
	Foo StringUnderlyingType = "underlying"
)

type Obj struct {
	Foo int    `json:"foo"`
	Bar string `json:"bar"`
}

type EntityWithAllPBTypes struct {
	Text               string               `orm:"text,omitempty"`
	TextUnderlyingType StringUnderlyingType `orm:"text_underlying,omitempty"`

	NumberFloat32 float32 `orm:"number_float32,omitempty"`
	NumberFloat64 float64 `orm:"number_float64,omitempty"`
	NumberInt     int     `orm:"number_int,omitempty"`
	NumberInt8    int8    `orm:"number_int8,omitempty"`
	NumberInt16   int16   `orm:"number_int16,omitempty"`
	NumberInt32   int32   `orm:"number_int32,omitempty"`
	NumberInt64   int64   `orm:"number_int64,omitempty"`
	NumberUint    uint    `orm:"number_uint,omitempty"`
	NumberUint8   uint8   `orm:"number_uint8,omitempty"`
	NumberUint16  uint16  `orm:"number_uint16,omitempty"`
	NumberUint32  uint32  `orm:"number_uint32,omitempty"`
	NumberUint64  uint64  `orm:"number_uint64,omitempty"`

	Bool   bool       `orm:"bool,omitempty"`
	Email  string     `orm:"email,omitempty"`
	Url    string     `orm:"url,omitempty"`
	Editor string     `orm:"editor,omitempty"`
	Date   *time.Time `orm:"date,omitempty"`

	JsonAsStringArray []string `orm:"json_string_array,omitempty"`
	JsonAsObject      Obj      `orm:"json_object,omitempty"`
	JsonAsObjectArray []Obj    `orm:"json_object_array,omitempty"`

	SingleRelation   string   `orm:"single_relation,omitempty"`
	MultipleRelation []string `orm:"multiple_relation,omitempty"`

	SingleSelect   string   `orm:"single_select,omitempty"`
	MultipleSelect []string `orm:"multiple_select,omitempty"`
}

func (_ EntityWithAllPBTypes) CollectionName() string {
	return "foo"

}

// Collection - for testing purpose only, note that this method is not part of Entity interface.
func (_ EntityWithAllPBTypes) Collection() *models.Collection {
	_schema := schema.NewSchema(
		&schema.SchemaField{Name: "text", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "text_underlying", Type: schema.FieldTypeText},

		&schema.SchemaField{Name: "number_float32", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_float64", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_int", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_int8", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_int16", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_int32", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_int64", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_uint", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_uint8", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_uint16", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_uint32", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "number_uint64", Type: schema.FieldTypeNumber},

		&schema.SchemaField{Name: "bool", Type: schema.FieldTypeBool},
		&schema.SchemaField{Name: "email", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "url", Type: schema.FieldTypeUrl},
		&schema.SchemaField{Name: "editor", Type: schema.FieldTypeEditor},
		&schema.SchemaField{Name: "date", Type: schema.FieldTypeDate},

		&schema.SchemaField{Name: "json_string_array", Type: schema.FieldTypeJson},
		&schema.SchemaField{Name: "json_object", Type: schema.FieldTypeJson},
		&schema.SchemaField{Name: "json_object_array", Type: schema.FieldTypeJson},

		&schema.SchemaField{Name: "single_relation", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MinSelect: pointer(0), MaxSelect: pointer(1)}},
		&schema.SchemaField{Name: "multiple_relation", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MinSelect: pointer(0), MaxSelect: pointer(10)}},

		&schema.SchemaField{Name: "single_select", Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
		&schema.SchemaField{Name: "multiple_select", Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
	)

	return &models.Collection{Name: "foo", Schema: _schema}
}

// hydrated in init function below
var (
	recordExample *models.Record
	entityExample EntityWithAllPBTypes
)

func init() {
	_date, err := time.Parse(types.DefaultDateLayout, "2023-05-12 19:51:05.000Z")
	if err != nil {
		panic(fmt.Sprintf("could not parse date: %v", err))
	}

	entityExample = EntityWithAllPBTypes{
		Text:               "foo",
		TextUnderlyingType: Foo,

		NumberFloat32: 32.25,
		NumberFloat64: 64.25,
		NumberInt:     -5,
		NumberInt8:    -8,
		NumberInt16:   -16,
		NumberInt32:   -32,
		NumberInt64:   -64,
		NumberUint:    5,
		NumberUint8:   8,
		NumberUint16:  16,
		NumberUint32:  32,
		NumberUint64:  64,

		Bool:   true,
		Email:  "foo@bar.qux",
		Url:    "http://foo.bar",
		Editor: "# Foo\nBar",
		Date:   &_date,

		JsonAsStringArray: []string{"foo", "bar"},
		JsonAsObject:      Obj{Foo: 12, Bar: "qux"},
		JsonAsObjectArray: []Obj{{Foo: 1, Bar: "1"}, {Foo: 2, Bar: "2"}},

		SingleRelation:   "foo",
		MultipleRelation: []string{"foo", "qux"},

		SingleSelect:   "foo",
		MultipleSelect: []string{"bar", "baz"},
	}

	recordExample = models.NewRecord(EntityWithAllPBTypes{}.Collection())
	recordExample.Set("text", "foo")
	recordExample.Set("text_underlying", "underlying")

	recordExample.Set("number_float32", 32.25)
	recordExample.Set("number_float64", 64.25)
	recordExample.Set("number_int", -5)
	recordExample.Set("number_int8", -8)
	recordExample.Set("number_int16", -16)
	recordExample.Set("number_int32", -32)
	recordExample.Set("number_int64", -64)
	recordExample.Set("number_uint", 5)
	recordExample.Set("number_uint8", 8)
	recordExample.Set("number_uint16", 16)
	recordExample.Set("number_uint32", 32)
	recordExample.Set("number_uint64", 64)

	recordExample.Set("bool", true)
	recordExample.Set("email", "foo@bar.qux")
	recordExample.Set("url", "http://foo.bar")
	recordExample.Set("editor", "# Foo\nBar")
	recordExample.Set("date", "2023-05-12 19:51:05.000Z")

	recordExample.Set("json_string_array", `["foo","bar"]`)
	recordExample.Set("json_object", `{"foo":12,"bar":"qux"}`)
	recordExample.Set("json_object_array", `[{"foo":1,"bar":"1"},{"foo":2,"bar":"2"}]`)

	recordExample.Set("single_relation", "foo")
	recordExample.Set("multiple_relation", `["foo","qux"]`)

	recordExample.Set("single_select", "foo")
	recordExample.Set("multiple_select", `["bar","baz"]`)
}
