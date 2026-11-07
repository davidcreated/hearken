package hearken

import (
	"reflect"
	"testing"
)

type expectedTarget struct {
	Device     string   `json:"device,omitempty"`
	Space      string   `json:"space,omitempty"`
	Devices    []string `json:"devices,omitempty"`
	Mode       string   `json:"mode,omitempty"`
	Ambiguous  []string `json:"ambiguous,omitempty"`
	Everything []string `json:"everything,omitempty"`
	None       bool     `json:"none,omitempty"`
}

func describe(t Target) expectedTarget {
	switch t := t.(type) {
	case DeviceTarget:
		return expectedTarget{Device: t.Device.ID}
	case SpaceTarget:
		return expectedTarget{Space: t.Space.ID, Devices: deviceIDs(t.Devices)}
	case ModeTarget:
		return expectedTarget{Mode: t.Mode.ID}
	case EverythingTarget:
		return expectedTarget{Everything: deviceIDs(t.Devices)}
	case Ambiguous:
		var ids []string
		for _, option := range t.Options {
			d := describe(option)
			ids = append(ids, d.Device+d.Space+d.Mode)
		}
		return expectedTarget{Ambiguous: ids}
	default:
		return expectedTarget{None: true}
	}
}

func deviceIDs(devices []Device) []string {
	var ids []string
	for _, d := range devices {
		ids = append(ids, d.ID)
	}
	return ids
}

func TestMatchSpec(t *testing.T) {
	var spec struct {
		Homes map[string]Home `json:"homes"`
		Cases []struct {
			Home      string         `json:"home"`
			Say       string         `json:"say"`
			ModesOnly bool           `json:"modesOnly"`
			Command   bool           `json:"command"`
			Expect    expectedTarget `json:"expect"`
		} `json:"cases"`
	}
	readSpec(t, "match.json", &spec)
	var parser Parser
	for _, c := range spec.Cases {
		home, ok := spec.Homes[c.Home]
		if !ok {
			t.Fatalf("%q: no home %q in the spec", c.Say, c.Home)
		}
		var got Target
		if c.Command {
			command, ok := parser.Parse(c.Say)
			if !ok {
				t.Errorf("%q: not understood as a command", c.Say)
				continue
			}
			got = Resolve(command, home)
		} else {
			got = Match(c.Say, home, c.ModesOnly)
		}
		if d := describe(got); !reflect.DeepEqual(d, c.Expect) {
			t.Errorf("%s / %q: got %+v, want %+v", c.Home, c.Say, d, c.Expect)
		}
	}
}
