package orm

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func DecodeAll[T Entity](records []*models.Record, entities []*T) error {
	if len(records) != len(entities) {
		return fmt.Errorf("length mismatch between records and entities provided")
	}

	for i, record := range records {
		if entities[i] == nil {
			var zeroValue T
			entities[i] = &zeroValue
		}

		if err := Decode(record, entities[i]); err != nil {
			return fmt.Errorf("could not decode %d element: %w", i, err)
		}
	}

	return nil
}

func Decode[T Entity](record *models.Record, entity *T) error {
	if record == nil {
		return fmt.Errorf("could not decode nil record")
	}

	if record.Collection() == nil {
		return fmt.Errorf("record's collection is nil")
	}

	collSchema := record.Collection().Schema

	if entity == nil {
		var zeroValue T
		entity = &zeroValue
	}

	ps := reflect.ValueOf(entity)
	s := ps.Elem()

	if s.Kind() != reflect.Struct {
		return fmt.Errorf("entity given is not a structure")
	}

	recordMap := RecordsColumnValueMap(record)

	numField := reflect.TypeOf(entity).Elem().NumField()
	for i := 0; i < numField; i++ {
		field := reflect.TypeOf(entity).Elem().Field(i)
		columnName := extractOrmNameFromTag(string(field.Tag))
		rawValue, ok := recordMap[columnName]
		if !ok {
			continue
		}

		fieldType := fieldFromColumnName(collSchema, columnName)
		if fieldType == nil {
			continue
		}

		entityField := reflect.ValueOf(entity).Elem().Field(i)

		switch fieldType.Type {
		case schema.FieldTypeText, schema.FieldTypeEmail, schema.FieldTypeUrl, schema.FieldTypeEditor:
			if strVal, ok := rawValue.(string); ok {
				entityField.SetString(strVal)
			}
			break

		case schema.FieldTypeNumber:
			if f64Val, ok := rawValue.(float64); ok {
				switch entityField.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					entityField.SetInt(int64(f64Val))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					entityField.SetUint(uint64(f64Val))
				case reflect.Float32, reflect.Float64:
					entityField.SetFloat(f64Val)
				}
			}
			break

		case schema.FieldTypeBool:
			if boolVal, ok := rawValue.(bool); ok {
				entityField.SetBool(boolVal)
			}
			break

		case schema.FieldTypeDate:
			datetime, err := time.Parse(types.DefaultDateLayout, fmt.Sprint(rawValue))
			if err != nil {
				break
			}
			entityField.Set(reflect.ValueOf(&datetime))
			break

		case schema.FieldTypeJson:
			bytesVal := []byte(fmt.Sprint(rawValue))

			val := reflect.New(entityField.Type()).Interface()
			if err := json.Unmarshal(bytesVal, val); err != nil {
				log.Printf("could not unmarshal %s: %v", string(bytesVal), err)
				break
			}
			entityField.Set(reflect.ValueOf(val).Elem())
			break

		case schema.FieldTypeRelation:
			if !(fieldType.Options.(*schema.RelationOptions)).IsMultiple() {
				if strVal, ok := rawValue.(string); ok {
					entityField.SetString(strVal)
				}
				break
			}

			rawJsonArray, ok := rawValue.(types.JsonArray[string])
			if !ok {
				log.Printf("could not cast %v to types.JsonArray[string]", rawValue)
				break
			}

			data, err := rawJsonArray.MarshalJSON()
			if err != nil {
				log.Printf("could not marshal %v", rawJsonArray)
				break
			}

			strSlice := []string{}
			if err := json.Unmarshal(data, &strSlice); err != nil {
				log.Printf("could not unmarshal %v", string(data))
				break
			}

			entityField.Set(reflect.ValueOf(strSlice))
			break

		case schema.FieldTypeSelect:
			if !(fieldType.Options.(*schema.SelectOptions)).IsMultiple() {
				if strVal, ok := rawValue.(string); ok {
					entityField.SetString(strVal)
				}
				break
			}

			rawJsonArray, ok := rawValue.(types.JsonArray[string])
			if !ok {
				log.Printf("could not cast %v to types.JsonArray[string]", rawValue)
				break
			}

			data, err := rawJsonArray.MarshalJSON()
			if err != nil {
				log.Printf("could not marshal %v", rawJsonArray)
				break
			}

			strSlice := []string{}
			if err := json.Unmarshal(data, &strSlice); err != nil {
				log.Printf("could not unmarshal %v", string(data))
				break
			}

			entityField.Set(reflect.ValueOf(strSlice))
			break
		}
	}

	return nil
}
