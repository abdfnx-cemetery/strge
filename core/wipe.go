package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gepis/strge"
	"github.com/gepis/strge/pkg/mflag"
)

func wipe(flags *mflag.FlagSet, action string, m storage.Store, args []string) int {
	err := m.Wipe()
	if jsonOutput {
		if err == nil {
			json.NewEncoder(os.Stdout).Encode(string(""))
		} else {
			json.NewEncoder(os.Stdout).Encode(err)
		}
	} else {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %+v\n", action, err)
		}
	}
	if err != nil {
		return 1
	}
	return 0
}

func init() {
	commands = append(commands, command{
		names:   []string{"wipe"},
		usage:   "Wipe all layers, images, and containers",
		minArgs: 0,
		action:  wipe,
		addFlags: func(flags *mflag.FlagSet, cmd *command) {
			flags.BoolVar(&jsonOutput, []string{"-json", "j"}, jsonOutput, "Prefer JSON output")
		},
	})
}
