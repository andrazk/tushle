package toshl

import (
	"fmt"

	"github.com/pkg/errors"
)

// Login with Toshl access token.
func (t *Toshl) Login() error {
	token, err := t.token()
	if err != nil {
		return errors.WithStack(err)
	}
	if token != "" {
		return nil
	}

	fmt.Fprintln(t.cli.Out(), "Go to https://developer.toshl.com/apps/ and create new personal token. Token will be securely stored.")
	fmt.Fprint(t.cli.Out(), "Personal token:")

	token, err = t.cli.ReadPassword()
	if err != nil {
		return errors.WithStack(err)
	}

	err = t.creds.Store(acronym, token)
	return errors.WithStack(err)
}
