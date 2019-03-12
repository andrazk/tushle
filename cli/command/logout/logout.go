package logout

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"tushle"
	"tushle/cli/banks"
)

// NewLogoutCommand creates a new `tushle logout` command
func NewLogoutCommand(tushleCli tushle.Cli, banks *banks.Repository) *cobra.Command {
	return &cobra.Command{
		Use:   "logout [bank]",
		Short: "Logout from your bank account. Toshl is default bank.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var b string
			if len(args) > 0 {
				b = args[0]
			}

			bank, err := banks.Available(b)
			if err != nil {
				return errors.WithStack(err)
			}

			return errors.WithStack(bank.Logout())
		},
	}
}
