package procra

import (
	"context"
	"github.com/Sirupsen/logrus"
)

type Proxy struct {
	*Bot
	Logger        *logrus.Logger
	URL           string
	Active        bool
	PrevState     ProxyState
	CurrentState  ProxyState
	CurrentConext context.Context
	fetchCh       chan *TargetURL
	activeCh      chan bool
	nextStateCh   chan ProxyState
	stopCh        chan bool
}

func NewProxy(url string, logger *logrus.Logger, state ProxyState) *Proxy {
	basecontext := context.Background()
	return &Proxy{
		URL:           url,
		Logger:        logger,
		Active:        true,
		CurrentState:  state,
		CurrentConext: basecontext,
		fetchCh:       make(chan *TargetURL),
		activeCh:      make(chan bool),
		nextStateCh:   make(chan ProxyState),
		stopCh:        make(chan bool),
	}
}

func (px *Proxy) Start() {
	for {
		select {
		case targurl := <-px.fetchCh:
			px.Logger.Info("fetch in px select loop")
			res, err := px.CurrentState.Fetch(targurl)
			ret := &Result{res, err}
			px.CurrentState.HandleResult(px, ret)
			px.Logger.Info("send result to bot")
			px.Logger.Info("bot", px.Bot)
			px.Bot.SendResult(ret)
			px.UpdateContext(ret)
		case next := <-px.nextStateCh:
			px.UpdateState(next)
		case active := <-px.activeCh:
			px.Active = active
		case <-px.stopCh:
			return
		}
	}
}

func (px *Proxy) Fetch(targ *TargetURL) {
	go func() {
		px.fetchCh <- targ
	}()
}

func (px *Proxy) Activate(active bool) {
	go func() {
		px.activeCh <- active
	}()
}

func (px *Proxy) NextState(state ProxyState) {
	go func() {
		px.nextStateCh <- state
	}()
}

func (px *Proxy) Stop() {
	go func() {
		px.stopCh <- true
	}()
}

func (px *Proxy) UpdateState(next ProxyState) {
	px.PrevState = px.CurrentState
	px.CurrentState = next
}

func (px *Proxy) UpdateContext(ret *Result) {

}
