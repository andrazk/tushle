package nkbm

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func (n *nkbm) Login() error {
	lc, err := n.getCreds()
	if err != nil {
		return errors.WithStack(err)
	}
	if lc != nil {
		return nil
	}

	lc = &localCred{
		Pass: randPass(),
	}

	if err := n.browser.Init(); err != nil {
		return errors.WithStack(err)
	}

	err = n.browser.Run(chromedp.Tasks{
		chromedp.Navigate(`https://bankanet.nkbm.si/prijava/bnk`),
		n.clearLogin(),
		n.enterUserPass(),
		n.enterOTPLogin(),
		chromedp.Click("#registerbrowser", chromedp.ByID),
		chromedp.Click("#logingumb > input.submit"),
		n.enterOTPSetup(),
		n.enterPassword(lc.Pass),

		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			var res []byte
			if err := n.getLocal33(&res).Do(ctx, h); err != nil {
				return errors.WithStack(err)
			}

			if len(res) == 0 {
				return errors.New("Credentials missing from browser local storage")
			}

			fmt.Println(string(res))

			lc.Store = string(res)
			err := n.creds.Store(acronym, *lc)
			return errors.WithStack(err)
		}),
	})

	return errors.WithStack(err)
}

func (n *nkbm) clearLogin() chromedp.Action {
	var res []byte
	return chromedp.Tasks{
		chromedp.Sleep(300 * time.Millisecond),
		chromedp.Evaluate(`window.localStorage.removeItem("_33");`, &res),
		chromedp.Reload(),
	}
}

func (n *nkbm) enterUserPass() chromedp.Action {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#indexForm`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			fmt.Fprintln(n.cli.Out(), "Login process for BankaNet started. You will be asked for username, password and OTP.")
			fmt.Fprint(n.cli.Out(), "Username:")
			username, err := n.cli.Read()
			if err != nil {
				return errors.WithStack(err)
			}

			fmt.Fprint(n.cli.Out(), "Password:")
			password, err := n.cli.ReadPassword()
			if err != nil {
				return errors.WithStack(err)
			}

			err = chromedp.Tasks{
				chromedp.SetValue("#p1", username, chromedp.ByID),
				chromedp.SetValue("#p2", password, chromedp.ByID),
			}.Do(ctx, h)
			return errors.WithStack(err)
		}),
		chromedp.Click("#logingumb > input"),
	}
}

func (n *nkbm) enterOTPSetup() chromedp.Action {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			fmt.Fprint(n.cli.Out(), "OTP Start Setup:")
			otp, err := n.cli.ReadPassword()
			if err != nil {
				return errors.WithStack(err)
			}

			err = chromedp.SetValue(`body > div.oobContainer > table.template1.requestOobContainer > tbody > tr:nth-child(3) > td.tlb > input[type="password"]`, otp).Do(ctx, h)
			return errors.WithStack(err)
		}),
		chromedp.Click("#bnkSplosniPogoji", chromedp.ByID),
		chromedp.Click("#okOobBtn", chromedp.ByID),
	}
}

func (n *nkbm) enterOTPLogin() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
		fmt.Fprint(n.cli.Out(), "OTP Login:")
		otp, err := n.cli.ReadPassword()
		if err != nil {
			return errors.WithStack(err)
		}

		err = chromedp.SetValue("#p2", otp, chromedp.ByID).Do(ctx, h)
		return errors.WithStack(err)
	})
}

func (n *nkbm) enterPassword(pass string) chromedp.Action {
	return chromedp.Tasks{
		chromedp.SetValue("body > div.nkbmcontent > center > div > table.template1 > tbody > tr:nth-child(3) > td.tlb > input", "tushle"),
		chromedp.SetValue(`body > div.nkbmcontent > center > div > table.template1 > tbody > tr:nth-child(4) > td.tlb > input`, pass),
		chromedp.SetValue(`body > div.nkbmcontent > center > div > table.template1 > tbody > tr:nth-child(5) > td.tlb > input`, pass),
		chromedp.Click("#okBtn", chromedp.ByID),
		chromedp.Sleep(300 * time.Millisecond),
		n.checkLoginError(),
	}
}

func (n *nkbm) getLocal33(res *[]byte) chromedp.Action {
	// Try 5 times
	tries := 5
	return chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
		for i := 0; i < tries; i++ {
			err := chromedp.Evaluate(`window.localStorage.getItem("_33");`, res).Do(ctx, h)
			if err != nil {
				return errors.WithStack(err)
			}
			if len(*res) > 0 {
				return nil
			}
			time.Sleep(300 * time.Millisecond)
		}
		return errors.Errorf("Credentials missing from browser local storage after %d tries", tries)
	})
}

func (n *nkbm) checkLoginError() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
		var errMsg string
		if err := chromedp.Text("#opozorila", &errMsg, chromedp.ByID).Do(ctx, h); err != nil {
			return errors.WithStack(err)
		}

		errMsg = strings.TrimSpace(errMsg)
		if errMsg != "" {
			return errors.New(errMsg)
		}

		return nil
	})
}

func (n *nkbm) getCreds() (*localCred, error) {
	var lc localCred
	if err := n.creds.Get(acronym, &lc); err != nil {
		return nil, errors.WithStack(err)
	}
	if lc.Pass == "" || lc.Store == "" {
		return nil, nil
	}
	return &lc, nil
}

var lowerRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var upperRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")
var specialRunes = []rune("!@#$%^&*()_")

func randPass() string {
	b := make([]rune, 16)
	for i := 0; i < 4; i++ {
		b[i] = lowerRunes[rand.Intn(len(lowerRunes))]
		b[i+4] = upperRunes[rand.Intn(len(upperRunes))]
		b[i+8] = numberRunes[rand.Intn(len(numberRunes))]
		b[i+12] = specialRunes[rand.Intn(len(specialRunes))]
	}
	return string(b)
}
