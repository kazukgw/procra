package procra

import (
	"context"

	"github.com/Sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

type Bot struct {
	Conf Config
	HTMLRepository
	Logger       *logrus.Logger
	BaseContext  context.Context
	PrevState    BotState
	CurrentState BotState
	Proxies      []*Proxy
	crn          *cron.Cron
	resultCh     chan *Result
	nextStateCh  chan BotState
	stopCh       chan error
	errCh        chan error
}

type Config struct {
	Cron string
}

func NewBot(conf Config, logger *logrus.Logger) *Bot {
	crn := cron.New()
	bot := &Bot{
		Conf: conf,
		crn:  crn,
	}
	return bot
}

func (bot *Bot) Start() error {
	bot.Logger.Info("bot start")
	go bot.crn.Start()

	for {
		select {
		case ret := <-bot.resultCh:
			// update context
			if err := bot.CurrentState.HandleResult(bot, ret); err != nil {
				bot.SendError(err)
				continue
			}
		case next := <-bot.nextStateCh:
			bot.UpdateState(next)
		case err := <-bot.errCh:
			if err := bot.CurrentState.HandleError(bot, err); err != nil {
				bot.Stop(err)
			}
		case err := <-bot.stopCh:
			return err
		}
	}
	return nil
}

func (bot *Bot) SendResult(result *Result) {
	go func() {
		bot.resultCh <- result
	}()
}

func (bot *Bot) SendNextState(state BotState) {
	go func() {
		bot.nextStateCh <- state
	}()
}

func (bot *Bot) SendError(err error) {
	go func() {
		bot.errCh <- err
	}()
}

func (bot *Bot) Stop(err error) {
	go func() {
		bot.stopCh <- err
	}()
}

func (bot *Bot) UpdateState(next BotState) {
	bot.crn.Stop()
	bot.crn = nil
	bot.PrevState = bot.CurrentState
	bot.CurrentState = next
	bot.crn = cron.New()
	bot.crn.AddFunc(bot.CurrentState.CronString(), func() {
		bot.CurrentState.Fetch(bot)
	})
	bot.crn.Start()
}
