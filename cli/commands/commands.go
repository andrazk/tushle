package commands

import (
	"github.com/spf13/cobra"

	"tushle"
	"tushle/cli/banks"
	"tushle/cli/command"
	"tushle/cli/command/login"
	"tushle/cli/command/logout"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, tushleCli tushle.Cli, banks *banks.Repository) {
	cmd.AddCommand(login.NewLoginCommand(tushleCli, banks))
	cmd.AddCommand(logout.NewLogoutCommand(tushleCli, banks))
	cmd.AddCommand(command.NewAccountsCommand(tushleCli, banks))
}
