package browser

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
)

// OptionFunc allows setting options on browser
type OptionFunc func(context.Context) chromedp.Option

// WithStandalone connects to existing Chrome instance on port 9222
func WithStandalone() OptionFunc {
	return func(ctx context.Context) chromedp.Option {
		return chromedp.WithTargets(client.New().WatchPageTargets(ctx))
	}
}

// WithLog sets log printer
func WithLog(f func(string, ...interface{})) OptionFunc {
	return func(ctx context.Context) chromedp.Option {
		return chromedp.WithLog(f)
	}
}
