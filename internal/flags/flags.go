package flags

import (
	"github.com/spf13/pflag"
)

var (
	All = &pflag.Flag{
		Name:      "all",
		Shorthand: "",
	}

	Tag = &pflag.Flag{
		Name:      "tag",
		Shorthand: "",
	}
)
