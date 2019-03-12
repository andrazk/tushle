package credentials

import (
	"encoding/json"

	"github.com/pkg/errors"
	keyring "github.com/zalando/go-keyring"

	"tushle"
)

// NewCredentials returns interface instance.
func NewCredentials(service string) tushle.Credentials {
	return &creds{
		service: service,
	}
}

type creds struct {
	service string
}

func (c *creds) Get(key string, v interface{}) error {
	s, err := keyring.Get(c.service, key)
	if err == keyring.ErrNotFound {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}
	if s == "" {
		return nil
	}

	err = json.Unmarshal([]byte(s), v)
	return errors.WithStack(err)
}

// Store saves the given credentials in the file store.
func (c *creds) Store(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return errors.WithStack(err)
	}

	err = keyring.Set(c.service, key, string(b))
	return errors.WithStack(err)
}

func (c *creds) Delete(key string) error {
	return errors.WithStack(keyring.Delete(c.service, key))
}
