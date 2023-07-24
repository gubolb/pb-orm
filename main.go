package orm

import (
	"fmt"

	"github.com/pocketbase/pocketbase/daos"
)

var defaultDao *daos.Dao
var ErrNotInit = fmt.Errorf("orm is not init")

func Setup(dao *daos.Dao) {
	defaultDao = dao
}
