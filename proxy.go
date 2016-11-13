package procra

import "context"

type Proxy struct {
	*Bot
	Active        bool
	PrevState     ProxyState
	CurrentState  ProxyState
	CurrentConext context.Context
	fetchCh       chan *TargetURL
	activeCh      chan bool
	nextStateCh   chan ProxyState
}

func (px *Proxy) Start() {
	for {
		select {
		case targurl := <-px.fetchCh:
			res, err := px.CurrentState.Fetch(targurl)
			ret := &Result{res, err}
			px.CurrentState.HandleResult(px, ret)
			px.Bot.SendResult(ret)
			px.UpdateContext(ret)
		case next := <-px.nextStateCh:
			px.UpdateState(next)
		case active := <-px.activeCh:
			px.Active = active
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

func (px *Proxy) UpdateState(next ProxyState) {
	px.PrevState = px.CurrentState
	px.CurrentState = next
}

func (px *Proxy) UpdateContext(ret *Result) {

}