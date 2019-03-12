package nkbm

import (
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"

	"tushle"
	"tushle/banks/nkbm/scrapper"
)

// Accounts scrapes account info
func (n *nkbm) Accounts() (tushle.Accounts, error) {
	var (
		accounts             tushle.Accounts
		names, balances, ids []string
	)

	lc, err := n.getCreds()
	if err != nil {
		return accounts, errors.WithStack(err)
	}
	if lc == nil {
		return accounts, errors.New("Missing credentials. Login to Nkbm first")
	}

	err = n.browser.Run(chromedp.Tasks{
		scrapper.Elogin(lc.Pass, lc.Store),
		chromedp.WaitReady(`#dock0`, chromedp.ByID),
		chromedp.Evaluate(scrapper.GetText(`.tlb_strcn > a`), &names),
		chromedp.Evaluate(scrapper.GetText(`tbody > tr > .tlb_txt`), &ids),
		chromedp.Evaluate(scrapper.GetText(`.tlb_zns > a`), &balances),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if nl, bl := len(names), len(balances); nl != bl {
		return accounts, errors.Errorf("Expecting same number of names and balances. Got %d names ans %d balances", nl, bl)
	}

	for i := range names {
		balance := strings.Split(balances[i], " ")
		if len(balance) != 2 {
			return accounts, errors.Errorf("Account balance missing currency. %s", names[i])
		}

		bal, err := toFloat(balance[0])
		if err != nil {
			return accounts, errors.Errorf("Account %s balance is not float. %v", names[i], err)
		}

		name := strings.Trim(names[i], ", ")

		accounts = append(accounts, &tushle.Account{
			ID:       strings.ToLower(name),
			Name:     name,
			Balance:  bal,
			Currency: strings.Trim(balance[1], ", "),
		})
	}

	return accounts, nil
}
