package command

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// NewTushleCli returns a TushleCli instance with IO output and error streams set by in, out and err.
func NewTushleCli(stdIn io.ReadCloser, stdOut, stdErr io.Writer) *TushleCli {
	return &TushleCli{
		err: stdErr,
		in:  stdIn,
		out: stdOut,
	}
}

// TushleCli is an instance the tushle command line client.
type TushleCli struct {
	err io.Writer
	out io.Writer
	in  io.Reader
}

// Err returns the writer used for stderr.
func (cli *TushleCli) Err() io.Writer {
	return cli.err
}

// Out returns the writer used for stdout.
func (cli *TushleCli) Out() io.Writer {
	return cli.out
}

// In returns the writer used for stdin.
func (cli *TushleCli) In() io.Reader {
	return cli.in
}

func (cli *TushleCli) Read() (string, error) {
	var text string
	var err error

	reader := bufio.NewReader(cli.in)
	for text == "" {
		text, err = reader.ReadString('\n')
		if err != nil {
			return text, errors.WithStack(err)
		}
		text = strings.TrimSpace(text)
	}
	return text, nil
}

// ReadPassword from stdin.
func (cli *TushleCli) ReadPassword() (string, error) {
	bytePass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", errors.WithStack(err)
	}
	fmt.Fprint(cli.out, "\n")
	return strings.TrimSpace(string(bytePass)), nil
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOutput(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
