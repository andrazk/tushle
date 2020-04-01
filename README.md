# Tushle (alpha)

Tushle is a command line tool which strives to provide simple to use interface for working with personal finance tracker Toshl and some banks.

It is also an experiment in designing command line tools, web scrapping with ChromeDP and other technologies that might catch my attention.

## Requirements

Currently, it's only supported on MacOS. Feel free to contact me if you are interested in giving it a go on other platforms.

## Install

Checkout the repository and run make command.

```bash
git checkout https://github.com/andrazk/tushle.git
make binary-osx
mv dist /usr/local/bin/tushle
```

## Usage

```bash
$ tushle --help
Tushle. Toshl CLI tool.

Usage:
  tushle
  tushle [command]

Available Commands:
  accounts    List accounts. Toshl is default bank.
  help        Help about the command
  login       Login to your bank account. Toshl is default bank.
  logout      Logout from your bank account. Toshl is default bank.

Flags:
      --config string      config file (default is $HOME/.tushle.yaml)
  -l, --log-level string   Set the logging level ("debug"|"info"|"warn"|"error"|"fatal") (default "info")
      --version            version for tushle

Use "tushle [command] --help" for more information about a command.
```

## CI

TODO

## License

Tushle is released under the terms of MIT License.
