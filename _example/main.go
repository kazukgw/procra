package main

import (
	"gopkg.in/h2non/gentleman.v1"
	"net/http"

	"bitbucket.org/monotaro/monotaro_crawler2"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open(
		"mysql",
		"monocra2:password@tcp(192.168.99.100:13306)/monocra2",
	)

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	logger := logrus.New()
	bot := monocra2.NewBot(
		logger,
		initBotState(db, logger),
		initProxies(db, logger),
	)
	bot.Start()
}

func initBotState(db *gorm.DB, logger *logrus.Logger) monocra2.BotState {
	return &DefaultBotState{
		DB:     db,
		Logger: logger,
	}
}

func initProxies(db *gorm.DB, logger *logrus.Logger) []*monocra2.Proxy {
	pxs := []*monocra2.Proxy{}
	urls := []string{"proxy1:8080", "proxy2:8080"}
	for _, u := range urls {
		st := &DefaultProxyState{
			DB:     db,
			Logger: logger,
		}
		px := monocra2.NewProxy(u, logger, st)
		pxs = append(pxs, px)
	}
	return pxs
}

type DefaultBotState struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

func (bs *DefaultBotState) String() string {
	return "default"
}

func (bs *DefaultBotState) CronString() string {
	return "@every 5s"
}

func (bs *DefaultBotState) Fetch(bot *monocra2.Bot) {
	bs.Logger.Info("fetch in bot state")
	px := bot.Proxies[0]
	targ := &monocra2.TargetURL{}
	targ.Scheme = "https"
	targ.Host = "google.co.jp"
	px.Fetch(targ)
}

func (bs *DefaultBotState) HandleResult(
	bot *monocra2.Bot,
	ret *monocra2.Result,
) error {
	bs.Logger.Info("handle result")
	return nil
}

func (bs *DefaultBotState) HandleError(
	bot *monocra2.Bot,
	err error,
) error {
	bs.Logger.Info("handle error")
	return nil
}

type DefaultProxyState struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

func (ps *DefaultProxyState) String() string {
	return "default"
}

func (ps *DefaultProxyState) Fetch(
	targ *monocra2.TargetURL,
) (*http.Response, error) {
	ps.Logger.Info("fetch in proxy")
	res, err := gentleman.New().URL(targ.URL().String()).Get().Do()
	ps.Logger.Info("res", res)
	return res.RawResponse, err
}

func (ps *DefaultProxyState) HandleResult(px *monocra2.Proxy, ret *monocra2.Result) error {
	ps.Logger.Info("handle result")
	ps.Logger.Info("result:", ret)
	return nil
}
