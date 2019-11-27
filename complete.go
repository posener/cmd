package cmd

import (
	"flag"

	"github.com/posener/complete/v2"
)

type completer SubCmd

func (c *completer) SubCmdList() []string {
	cmd := (*SubCmd)(c)
	return cmd.subNames()
}

func (c *completer) SubCmdGet(name string) complete.Completer {
	if c.sub[name] == nil {
		return nil
	}
	return (*completer)(c.sub[name])
}

func (c *completer) FlagList() []string {
	if len(c.sub) != 0 {
		return nil
	}
	var flags []string
	c.flagSet.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f.Name)
	})
	return flags
}

func (c *completer) FlagGet(flag string) complete.Predictor {
	f := c.flagSet.Lookup(flag)
	if f == nil {
		return nil
	}
	if p, ok := f.Value.(complete.Predictor); ok {
		return p
	}
	return nil
}

func (c *completer) ArgsGet() complete.Predictor {
	if c.args.predict.Predictor != nil {
		return c.args.predict
	}
	if comp, ok := c.args.value.(complete.Predictor); ok {
		return comp
	}
	return nil
}
