package procra

import (
	"github.com/jinzhu/gorm"
	"net/url"
)

type TargetURL struct {
	gorm.Model

	Scheme   string `gorm:"index:idx_scheme"`
	Host     string `gorm:"index:idx_host"`
	Path     string `gorm:"index:idx_path"`
	RawQuery string `gorm:"index:idx_rawquery"`
	Fragment string `gorm:"index:idx_fragment"`

	Rank     int `gorm:"index:idx_rank"`
	Priority int `gorm:"index:idx_priority"`

	target *url.URL `gorm:"-"`
}

func (targ *TargetURL) URL() *url.URL {
	if targ.target == nil {
		targ.target = &url.URL{
			Scheme:   targ.Scheme,
			Host:     targ.Host,
			Path:     targ.Path,
			RawQuery: targ.RawQuery,
			Fragment: targ.Fragment,
		}
	}
	return targ.target
}

func NextTargetURL(db *gorm.DB) *TargetURL {
	targ := &TargetURL{}
	db.Joins(`
	join target_url_stats on
	target_urls.id = target_url_stats.target_url_id
	`).Order(`
	target_url_stats.total_attempted_num asc,
	target_url_stats.last_attempted asc
	`).First(targ)
	return targ
}
