package monocra2

type BotState interface {
	String() string
	CronString() string
	Fetch(*Bot)
	HandleResult(*Bot, *Result) error
	HandleError(*Bot, error) error
}

type SleepState struct {
	SleepTime string
	NextState BotState
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
