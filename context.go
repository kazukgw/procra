package monocra2

import "context"

func NextURLFromCtx(ctx context.Context) *TargetURL {
	targurl, _ := ctx.Value("next").(*TargetURL)
	return targurl
}

func HistoryFromCtx(ctx context.Context) []*Result {
	hist, _ := ctx.Value("history").([]*Result)
	return hist
}

// func LastResultFromCtx(ctx context.Context) *Result {
// 	return
// }
