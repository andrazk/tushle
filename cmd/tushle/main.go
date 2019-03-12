package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"tushle/banks/nkbm"
	"tushle/banks/toshl"
	"tushle/cli"
	"tushle/cli/banks"
	"tushle/cli/command"
	"tushle/cli/commands"
	cliflags "tushle/cli/flags"
	"tushle/lib/browser"
	"tushle/lib/credentials"
)

func main() {
	// Set terminal emulation based on platform as required.
	stdin, stdout, stderr := stdStreams()
	logrus.SetOutput(stdout)

	tushleCli := command.NewTushleCli(stdin, stdout, stderr)

	// Credentials
	creds := credentials.NewCredentials("tushle")

	// Banks
	banks := banks.NewRepository()
	{
		// Toshl bank
		bank := toshl.New(tushleCli, creds)
		banks.AddDefaultBank(bank)
	}
	{
		// Nkbm bank
		browser, err := browser.New(browser.WithLog(logrus.Debugf))
		if err != nil {
			fmt.Fprintln(stderr, err)
			os.Exit(1)
		}
		// defer browser.Shutdown()
		bank := nkbm.NewNkbm(tushleCli, creds, browser)
		banks.AddBank(bank)
	}

	cmd := newTushleCommand(tushleCli, banks)

	if err := cmd.Execute(); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(stderr, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
}

func newTushleCommand(tushleCli *command.TushleCli, banks *banks.Repository) *cobra.Command {
	opts := cliflags.NewOptions()
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Use:              "tushle",
		Short:            "Tushle. Toshl CLI tool.",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.ShowHelp(tushleCli.Err())(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Version:               fmt.Sprintf("%s, build %s", cli.Version, cli.GitCommit),
		DisableFlagsInUseLine: true,
	}
	cli.SetupRootCommand(cmd)

	flags = cmd.PersistentFlags()
	flags.StringVar(&opts.ConfigFile, "config", "", "config file (default is $HOME/.tushle.yaml)")
	flags.StringVarP(&opts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)

	switch opts.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	}
	fmt.Println("LOG LEVEV:", opts.LogLevel)

	cmd.SetOutput(tushleCli.Out())
	commands.AddCommands(cmd, tushleCli, banks)

	return cmd
}

// stdStreams returns the standard streams (stdin, stdout, stderr).
func stdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"tushle: '%s' is not a command.\nSee 'tushle --help'", args[0])
}
