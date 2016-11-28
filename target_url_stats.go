package monocra2

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TargetURLStats struct {
	gorm.Model

	TotalSavedNum     int        `gorm:"not null;index:idx_total_saved_num" sql:"DEFAULT:0"`
	TotalAttemptedNum int        `gorm:"not null;index:idx_total_attemtped_num" sql:"DEFAULT:0"`
	LastAttempted     *time.Time `gorm:"index:idx_last_attemted"`
	LastResult        string     `gorm:"not null;index:idx_last_result"`

	TargetURLID uint `gorm:"not null;index:idx_target_url_id"`
}
