package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gepis/strge"
	"github.com/gepis/strge/internal/opt"
	"github.com/gepis/strge/pkg/mflag"
)

func getNames(flags *mflag.FlagSet, action string, m storage.Store, args []string) int {
	if len(args) < 1 {
		return 1
	}
	id, err := m.Lookup(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	names, err := m.Names(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	if jsonOutput {
		json.NewEncoder(os.Stdout).Encode(append([]string{}, names...))
	} else {
		for _, name := range names {
			fmt.Printf("%s\n", name)
		}
	}
	return 0
}

func addNames(flags *mflag.FlagSet, action string, m storage.Store, args []string) int {
	if len(args) < 1 {
		return 1
	}
	id, err := m.Lookup(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	oldnames, err := m.Names(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	newNames := []string{}
	if oldnames != nil {
		newNames = append(newNames, oldnames...)
	}
	if paramNames != nil {
		newNames = append(newNames, paramNames...)
	}
	if err := m.SetNames(id, newNames); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	names, err := m.Names(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	if jsonOutput {
		json.NewEncoder(os.Stdout).Encode(names)
	}
	return 0
}

func setNames(flags *mflag.FlagSet, action string, m storage.Store, args []string) int {
	if len(args) < 1 {
		return 1
	}
	id, err := m.Lookup(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	if err := m.SetNames(id, paramNames); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	names, err := m.Names(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return 1
	}
	if jsonOutput {
		json.NewEncoder(os.Stdout).Encode(names)
	}
	return 0
}

func init() {
	commands = append(commands, command{
		names:       []string{"get-names", "getnames"},
		optionsHelp: "[options [...]] imageOrContainerNameOrID",
		usage:       "Get layer, image, or container name or names",
		minArgs:     1,
		action:      getNames,
		addFlags: func(flags *mflag.FlagSet, cmd *command) {
			flags.BoolVar(&jsonOutput, []string{"-json", "j"}, jsonOutput, "Prefer JSON output")
		},
	})
	commands = append(commands, command{
		names:       []string{"add-names", "addnames"},
		optionsHelp: "[options [...]] imageOrContainerNameOrID",
		usage:       "Add layer, image, or container name or names",
		minArgs:     1,
		action:      addNames,
		addFlags: func(flags *mflag.FlagSet, cmd *command) {
			flags.Var(opt.NewListOptRef(&paramNames, nil), []string{"-name", "n"}, "New name")
			flags.BoolVar(&jsonOutput, []string{"-json", "j"}, jsonOutput, "Prefer JSON output")
		},
	})
	commands = append(commands, command{
		names:       []string{"set-names", "setnames"},
		optionsHelp: "[options [...]] imageOrContainerNameOrID",
		usage:       "Set layer, image, or container name or names",
		minArgs:     1,
		action:      setNames,
		addFlags: func(flags *mflag.FlagSet, cmd *command) {
			flags.Var(opt.NewListOptRef(&paramNames, nil), []string{"-name", "n"}, "New name")
			flags.BoolVar(&jsonOutput, []string{"-json", "j"}, jsonOutput, "Prefer JSON output")
		},
	})
}
