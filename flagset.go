package subcmd

import (
	"flag"
	"io"
	"time"
)

// flagSet wraps flag.FlagSet in order to hide APIs.
type flagSet struct {
	*flag.FlagSet
}

func (f *flagSet) Bool(name string, value bool, usage string) *bool {
	return f.FlagSet.Bool(name, value, usage)
}

func (f *flagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	return f.FlagSet.Duration(name, value, usage)
}

func (f *flagSet) Float64(name string, value float64, usage string) *float64 {
	return f.FlagSet.Float64(name, value, usage)
}

func (f *flagSet) Int(name string, value int, usage string) *int {
	return f.FlagSet.Int(name, value, usage)
}

func (f *flagSet) Int64(name string, value int64, usage string) *int64 {
	return f.FlagSet.Int64(name, value, usage)
}

func (f *flagSet) Output() io.Writer {
	return f.FlagSet.Output()
}

func (f *flagSet) Parsed() bool {
	return f.FlagSet.Parsed()
}

func (f *flagSet) Set(name, value string) error {
	return f.FlagSet.Set(name, value)
}

func (f *flagSet) String(name string, value string, usage string) *string {
	return f.FlagSet.String(name, value, usage)
}

func (f *flagSet) Uint(name string, value uint, usage string) *uint {
	return f.FlagSet.Uint(name, value, usage)
}

func (f *flagSet) Uint64(name string, value uint64, usage string) *uint64 {
	return f.FlagSet.Uint64(name, value, usage)
}

func (f *flagSet) Var(value flag.Value, name string, usage string) {
	f.FlagSet.Var(value, name, usage)
}

func (f *flagSet) Visit(fn func(*flag.Flag)) {
	f.FlagSet.Visit(fn)
}

func (f *flagSet) VisitAll(fn func(*flag.Flag)) {
	f.FlagSet.VisitAll(fn)
}
