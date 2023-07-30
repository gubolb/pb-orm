package orm

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func EncodeAll[T Entity](entities []*T, dao *daos.Dao) ([]*models.Record, error) {
	records := make([]*models.Record, len(entities))
	for i, e := range entities {
		record, err := Encode(e, dao)
		if err != nil {
			return nil, fmt.Errorf("could not encode %d element: %w", i, err)
		}
		records[i] = record
	}
	return records, nil
}

func Encode[T Entity](entity *T, dao *daos.Dao) (*models.Record, error) {
	if dao == nil {
		return nil, fmt.Errorf("could not encode: dao is nil")
	}

	if entity == nil {
		return nil, fmt.Errorf("could not encode nil entity")
	}

	coll, err := dao.FindCollectionByNameOrId((*entity).CollectionName())
	if err != nil {
		return nil, fmt.Errorf("could not get entity collection: %w", err)
	} else if coll == nil {
		return nil, fmt.Errorf("could not get entity collection: collection is nil")
	}

	ps := reflect.ValueOf(entity)
	s := ps.Elem()

	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity given is not a structure")
	}

	r := models.NewRecord(coll)

	numField := reflect.TypeOf(entity).Elem().NumField()
	for i := 0; i < numField; i++ {
		field := reflect.TypeOf(entity).Elem().Field(i)
		columnName := extractOrmNameFromTag(string(field.Tag))

		fieldType := fieldFromColumnName(coll.Schema, columnName)
		if fieldType == nil {
			continue
		}

		entityField := reflect.ValueOf(entity).Elem().Field(i)

		if isOmitable(string(field.Tag)) {
			if entityField.IsZero() {
				break
			}
		}

		switch fieldType.Type {
		case schema.FieldTypeText, schema.FieldTypeEmail, schema.FieldTypeUrl, schema.FieldTypeEditor:
			if field.Type.Kind() != reflect.String {
				break
			}
			r.Set(columnName, entityField.String())
			break

		case schema.FieldTypeNumber:
			switch entityField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				r.Set(columnName, entityField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				r.Set(columnName, entityField.Uint())
			case reflect.Float32, reflect.Float64:
				r.Set(columnName, entityField.Float())
			}
			break

		case schema.FieldTypeBool:
			if field.Type.Kind() != reflect.Bool {
				break
			}
			r.Set(columnName, entityField.Bool())
			break

		case schema.FieldTypeDate:
			_time, ok := entityField.Interface().(*time.Time)
			if !ok || _time == nil {
				break
			}

			r.Set(columnName, _time.String())
			break

		case schema.FieldTypeJson:
			data, err := json.Marshal(entityField.Interface())
			if err != nil {
				log.Printf("could not marshal %v: %v", entityField.Interface(), err)
				break
			}

			r.Set(columnName, string(data))
			break

		case schema.FieldTypeRelation:
			if !(fieldType.Options.(*schema.RelationOptions)).IsMultiple() {
				if field.Type.Kind() != reflect.String {
					break
				}
				r.Set(columnName, entityField.String())
				break
			}

			data, err := json.Marshal(entityField.Interface())
			if err != nil {
				log.Printf("could not marshal %v: %v", entityField.Interface(), err)
				break
			}

			r.Set(columnName, string(data))
			break

		case schema.FieldTypeSelect:
			if !(fieldType.Options.(*schema.SelectOptions)).IsMultiple() {
				if field.Type.Kind() != reflect.String {
					break
				}
				r.Set(columnName, entityField.String())
				break
			}

			data, err := json.Marshal(entityField.Interface())
			if err != nil {
				log.Printf("could not marshal %v: %v", entityField.Interface(), err)
				break
			}

			r.Set(columnName, string(data))
			break
		}
	}

	return r, nil
}
