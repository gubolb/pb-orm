package orm

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

var pbMetadata = []string{schema.FieldNameCreated, schema.FieldNameUpdated}

// RecordsColumnValueMap returns business column value map
// (i.e. r.ColumnValueMap without PocketBase metadata such as created and updated).
func RecordsColumnValueMap(r *models.Record) map[string]interface{} {
	if r == nil {
		return nil
	}

	cvm := r.ColumnValueMap()
	for _, metadata := range pbMetadata {
		delete(cvm, metadata)
	}
	return cvm
}
