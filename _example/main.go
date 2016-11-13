package main

import (
	"gopkg.in/h2non/gentleman.v1"
	"net/http"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kazukgw/procra"
)

func main() {
	db, err := gorm.Open(
		"mysql",
		"procra:password@tcp(192.168.99.100:13306)/procradb",
	)

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	logger := logrus.New()
	bot := procra.NewBot(
		logger,
		initBotState(db, logger),
		initProxies(db, logger),
	)
	bot.Start()
}

func initBotState(db *gorm.DB, logger *logrus.Logger) procra.BotState {
	return &DefaultBotState{
		DB:     db,
		Logger: logger,
	}
}

func initProxies(db *gorm.DB, logger *logrus.Logger) []*procra.Proxy {
	pxs := []*procra.Proxy{}
	urls := []string{"proxy1:8080", "proxy2:8080"}
	for _, u := range urls {
		st := &DefaultProxyState{
			DB:     db,
			Logger: logger,
		}
		px := procra.NewProxy(u, logger, st)
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

func (bs *DefaultBotState) Fetch(bot *procra.Bot) {
	bs.Logger.Info("fetch in bot state")
	px := bot.Proxies[0]
	targ := &procra.TargetURL{}
	targ.Scheme = "https"
	targ.Host = "google.co.jp"
	px.Fetch(targ)
}

func (bs *DefaultBotState) HandleResult(
	bot *procra.Bot,
	ret *procra.Result,
) error {
	bs.Logger.Info("handle result")
	return nil
}

func (bs *DefaultBotState) HandleError(
	bot *procra.Bot,
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
	targ *procra.TargetURL,
) (*http.Response, error) {
	ps.Logger.Info("fetch in proxy")
	res, err := gentleman.New().URL(targ.URL().String()).Get().Do()
	ps.Logger.Info("res", res)
	return res.RawResponse, err
}

func (ps *DefaultProxyState) HandleResult(px *procra.Proxy, ret *procra.Result) error {
	ps.Logger.Info("handle result")
	ps.Logger.Info("result:", ret)
	return nil
}
