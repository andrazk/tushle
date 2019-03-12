package toshl

import (
	"fmt"
	"net/http"

	"tushle"

	"github.com/pkg/errors"
)

// Accounts returns the list of users accounts
func (t *Toshl) Accounts() (tushle.Accounts, error) {
	var accounts tushle.Accounts
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts", baseURL), nil)
	if err != nil {
		return accounts, errors.WithStack(err)
	}

	var v []Account
	if err := t.do(req, &v); err != nil {
		return accounts, errors.WithStack(err)
	}

	for _, a := range v {
		accounts = append(accounts, &tushle.Account{
			ID:       a.ID,
			Balance:  a.Balance,
			Currency: a.Currency.Code,
			Name:     a.Name,
		})
	}
	return accounts, nil
}
