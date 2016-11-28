package monocra2

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

	Rank     int `gorm:"index:idx_rank" sql:"DEFAULT:1"`
	Priority int `gorm:"index:idx_priority" sql:"DEFAULT:1"`

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

func (targ *TargetURL) AfterCreate(tx *gorm.DB) {
	tx.Create(&TargetURLStats{TargetURLID: targ.ID})
}

func NewTargetURLFromRawURL(rawurl string) (*TargetURL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	targ := &TargetURL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
		target:   u,
	}
	return targ, nil
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

func NextTargetURLByRandom(db *gorm.DB) *TargetURL {
	return nil
}
