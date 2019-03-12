package nkbm

import "github.com/pkg/errors"

func (n *nkbm) Logout() error {
	return errors.WithStack(n.creds.Delete(acronym))
}
