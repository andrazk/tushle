package browser

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

// Browser interface
type Browser interface {
	Init() error
	Run(chromedp.Action) error
	Shutdown() error
}

// New Browser
func New(options ...OptionFunc) (Browser, error) {
	ctx, cancel := context.WithCancel(context.Background())

	opts := []chromedp.Option{}
	for _, of := range options {
		opts = append(opts, of(ctx))
	}

	return &cdp{
		ctx:    ctx,
		cancel: cancel,
		opts:   opts,
	}, nil
}

type cdp struct {
	ctx    context.Context
	cancel context.CancelFunc
	cdp    *chromedp.CDP
	opts   []chromedp.Option
}

func (c *cdp) Init() error {
	var err error
	if c.cdp != nil {
		return nil
	}

	c.cdp, err = chromedp.New(c.ctx, c.opts...)
	return errors.WithStack(err)
}

func (c *cdp) Run(a chromedp.Action) error {
	if c.cdp == nil {
		if err := c.Init(); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := c.cdp.Run(c.ctx, a); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *cdp) Shutdown() error {
	defer c.cancel()

	if c.cdp == nil {
		return nil
	}

	if err := c.cdp.Shutdown(c.ctx); err != nil {
		return errors.WithStack(err)
	}

	// wait for chrome to finish
	if err := c.cdp.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
