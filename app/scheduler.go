package app

import "context"

func (a *App) runScheduler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		}
	}
}
