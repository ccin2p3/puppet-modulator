package flagutils

import "github.com/spf13/pflag"

func HasFlag(set *pflag.FlagSet, key string) bool {
	var seen bool
	set.Visit(func(f *pflag.Flag) {
		if f.Name == key {
			seen = true
		}
	})
	return seen
}
