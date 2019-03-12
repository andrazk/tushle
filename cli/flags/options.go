package flags

// NewOptions returns a new Options.
func NewOptions() *Options {
	return &Options{}
}

// Options are the options used to configure the cli.
type Options struct {
	ConfigFile string
	Debug      bool
	LogLevel   string
}
