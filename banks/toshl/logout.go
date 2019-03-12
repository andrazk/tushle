package toshl

import "github.com/pkg/errors"

// Logout from toshl.
func (t *Toshl) Logout() error {
	return errors.WithStack(t.creds.Delete(acronym))
}
