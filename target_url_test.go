package procra

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func TestNextTargetURL(t *testing.T) {
	// db, err := gorm.Open("mysql", "root:password@tcp(192.168.99.100:13306)/procra")
	// db = db.Debug()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// targ := NextTargetURL(db)
	// if targ == nil {
	// 	t.Error()
	// }
}
