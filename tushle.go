package tushle

import "io"

// Banker describes interface for interaction with bank API.
type Banker interface {
	Accounts() (Accounts, error)
	Acronym() string
	Login() error
	Logout() error
}

// Cli represents the tushle command line client.
type Cli interface {
	Err() io.Writer
	Out() io.Writer
	In() io.Reader
	Read() (string, error)
	ReadPassword() (string, error)
}

// Credentials interface.
type Credentials interface {
	Get(key string, v interface{}) error
	Store(key string, v interface{}) error
	Delete(key string) error
}

// Account holds account data
type Account struct {
	ID       string
	Name     string
	Balance  float64
	Currency string
}

// Accounts holds list of accounts
type Accounts []*Account
