// Command hearken parses spoken commands and matches them to a home, for
// trying the library from a terminal.
//
//	hearken parse "turn off the living room heater"
//	hearken match --home examples/home.json "set the kitchen to away mode"
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/davidcreated/hearken"
)

const usage = `usage:
  hearken parse [--wake pattern] "what was said"
  hearken match --home home.json [--wake pattern] "what was said"`

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%s", usage)
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	wake := flags.String("wake", "", "wake phrase pattern to ignore (a regular expression)")
	homePath := flags.String("home", "", "home JSON file (devices, spaces, modes)")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("%s", usage)
	}
	parser, err := hearken.NewParser(*wake)
	if err != nil {
		return fmt.Errorf("--wake: %w", err)
	}
	command, ok := parser.Parse(flags.Arg(0))
	if !ok {
		return fmt.Errorf("not a command: %q", flags.Arg(0))
	}

	switch args[0] {
	case "parse":
		return write(out, command)
	case "match":
		if *homePath == "" {
			return fmt.Errorf("match needs --home\n%s", usage)
		}
		raw, err := os.ReadFile(*homePath)
		if err != nil {
			return err
		}
		var home hearken.Home
		if err := json.Unmarshal(raw, &home); err != nil {
			return fmt.Errorf("%s: %w", *homePath, err)
		}
		return write(out, map[string]any{
			"command": command,
			"target":  describe(hearken.Resolve(command, home)),
		})
	default:
		return fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
}

func describe(target hearken.Target) map[string]any {
	switch t := target.(type) {
	case hearken.DeviceTarget:
		return map[string]any{"device": t.Device}
	case hearken.SpaceTarget:
		return map[string]any{"space": t.Space, "devices": t.Devices}
	case hearken.ModeTarget:
		return map[string]any{"mode": t.Mode}
	case hearken.EverythingTarget:
		return map[string]any{"everything": t.Devices}
	case hearken.Ambiguous:
		options := make([]map[string]any, 0, len(t.Options))
		for _, option := range t.Options {
			options = append(options, describe(option))
		}
		return map[string]any{"ambiguous": options}
	default:
		return map[string]any{"none": true}
	}
}

func write(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
