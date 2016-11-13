package state

type SleepState struct {
	SleepTime string
	NextState BotState
}

func (ss SleepState) String() string {
	return "sleep"
}

func (ss SleepState) CronString() string {
	return "@every " + ss.SleepTime
}

func (ss SleepState) Fetch(bot *Bot) {
	bot.SendResult(&Result{})
}

func (ss SleepState) HandleResult(bot *Bot, ret *Result) error {
	bot.SendNextState(ss.NextState)
	return nil
}

func (ss SleepState) HandleError(bot *Bot, ret *Result) error {
	return nil
}
