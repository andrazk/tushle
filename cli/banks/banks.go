package banks

import (
	"tushle"

	"github.com/pkg/errors"
)

// NewRepository returns Repository instance.
func NewRepository() *Repository {
	return &Repository{
		available: map[string]tushle.Banker{},
	}
}

// Repository of available banks.
type Repository struct {
	available map[string]tushle.Banker
}

// AddBank to available banks.
func (r *Repository) AddBank(bank tushle.Banker) {
	r.available[bank.Acronym()] = bank
}

// AddDefaultBank to available banks.
func (r *Repository) AddDefaultBank(bank tushle.Banker) {
	r.available[bank.Acronym()] = bank
	r.available[""] = bank
}

// Available returns bank instance if available.
func (r *Repository) Available(key string) (tushle.Banker, error) {
	bank, ok := r.available[key]
	if !ok {
		return nil, errors.Errorf("Bank %s not available", key)
	}

	return bank, nil
}
