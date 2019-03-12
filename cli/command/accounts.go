package command

import (
	"fmt"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"tushle"
	"tushle/cli/banks"
)

// NewAccountsCommand creates a new `tushle accounts` command.
func NewAccountsCommand(tushleCli tushle.Cli, banks *banks.Repository) *cobra.Command {
	return &cobra.Command{
		Use:   "accounts [bank]",
		Short: "List accounts. Toshl is default bank.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var b string
			if len(args) > 0 {
				b = args[0]
			}

			bank, err := banks.Available(b)
			if err != nil {
				return errors.WithStack(err)
			}

			accounts, err := bank.Accounts()
			if err != nil {
				return errors.WithStack(err)
			}

			w := tabwriter.NewWriter(tushleCli.Out(), 3, 1, 3, ' ', tabwriter.Debug)
			fmt.Fprintln(w, "ID\tName\tBalance\t")
			for _, a := range accounts {
				fmt.Fprintf(w, "%s\t%s\t%.2f %s\t\n", a.ID, a.Name, a.Balance, a.Currency)
			}
			return errors.WithStack(w.Flush())
		},
	}
}
