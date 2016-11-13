package procra

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TargetURLStats struct {
	gorm.Model

	TotalSavedNum     int       `gorm:"index:idx_total_saved_num"`
	TotalAttemptedNum int       `gorm:"index:idx_total_attemtped_num"`
	LastAttempted     time.Time `gorm:"not null;index:idx_last_attemted"`
	LastResult        string    `gorm:"not null;index:idx_last_result""`

	TargetURLID int `gorm:"index:idx_target_url_id"`
}
