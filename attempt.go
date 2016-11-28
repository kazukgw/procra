package monocra2

import (
	"github.com/jinzhu/gorm"
)

type Attempt struct {
	gorm.Model

	Result      string
	StatusCode  int
	TargetURLID uint `gorm:"not null;index:idx_target_url_id"`
}
