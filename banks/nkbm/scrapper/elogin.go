package scrapper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

// Elogin using e-pass to Nkbm.
func Elogin(pass, localStorage string) chromedp.Action {
	selInput := `#p2`

	localStorage = strings.Trim(localStorage, `"`)

	return chromedp.Tasks{
		chromedp.Navigate(`https://bankanet.nkbm.si/prijava/bnk`),
		chromedp.Sleep(300 * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			var (
				exists bool
				res    []byte
			)

			if err := chromedp.Evaluate(`document.getElementById('indexForm') !== undefined`, &exists).Do(ctx, h); err != nil {
				return errors.WithStack(err)
			}
			if !exists {
				// Local storage already set
				return nil
			}

			return chromedp.Tasks{
				chromedp.Evaluate(fmt.Sprintf(`
					window.localStorage.setItem("_33", '%s');
					window.localStorage.setItem("authTypeB_33", '{"bnk":{"authType":"t"}}');
				`, localStorage), &res),
				chromedp.Reload(),
			}.Do(ctx, h)
		}),
		// chromedp.WaitReady(selInput, chromedp.ByID),
		chromedp.SetValue(selInput, pass, chromedp.ByID),
		chromedp.Click(`#loginBtn`, chromedp.ByID),
		chromedp.WaitVisible(`#dock0`),
	}
}
