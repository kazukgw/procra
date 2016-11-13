package procra

import (
	"github.com/jinzhu/gorm"
)

type Attempt struct {
	gorm.Model

	Result      string
	StatusCode  int
	TargetURLID int `gorm:"index:idx_target_url_id"`
}
