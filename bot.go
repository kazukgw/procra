package monocra2

import (
	"context"

	"github.com/Sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

type Bot struct {
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

func NewBot(
	logger *logrus.Logger,
	state BotState,
	pxs []*Proxy,
) *Bot {
	bot := &Bot{
		Logger:       logger,
		Proxies:      pxs,
		CurrentState: state,
		resultCh:     make(chan *Result),
		nextStateCh:  make(chan BotState),
		stopCh:       make(chan error),
		errCh:        make(chan error),
	}
	for _, px := range pxs {
		px.Bot = bot
	}
	return bot
}

func (bot *Bot) Start() error {
	bot.Logger.Info("bot start")
	bot.UpdateState(bot.CurrentState)
	bot.Logger.Info(bot.CurrentState)
	bot.StartProxies()
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
			bot.crn.Start()
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
		bot.Logger.Info("==> bot", bot)
		bot.Logger.Info("==> result", result)
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
	if bot.crn != nil {
		bot.crn.Stop()
		bot.crn = nil
	}
	bot.PrevState = bot.CurrentState
	bot.CurrentState = next
	bot.crn = cron.New()
	bot.crn.AddFunc(bot.CurrentState.CronString(), func() {
		bot.Logger.Debug("called cron fetch func")
		bot.CurrentState.Fetch(bot)
	})
}

func (bot *Bot) StartProxies() {
	for _, px := range bot.Proxies {
		go px.Start()
	}
}
