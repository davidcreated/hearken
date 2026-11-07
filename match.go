package hearken

import (
	"regexp"
	"strings"
)

// Target is what a spoken name was matched to: a DeviceTarget, SpaceTarget,
// ModeTarget, EverythingTarget, Ambiguous or NoMatch.
type Target interface{ target() }

// DeviceTarget is one device.
type DeviceTarget struct{ Device Device }

// SpaceTarget is a room and the devices in it.
type SpaceTarget struct {
	Space   Space
	Devices []Device
}

// ModeTarget is one mode.
type ModeTarget struct{ Mode Mode }

// EverythingTarget is every device in the home.
type EverythingTarget struct{ Devices []Device }

// Ambiguous means several things fit equally well; ask which one.
type Ambiguous struct {
	Spoken  string
	Options []Target
}

// NoMatch means nothing in the home has that name.
type NoMatch struct{ Spoken string }

func (DeviceTarget) target()     {}
func (SpaceTarget) target()      {}
func (ModeTarget) target()       {}
func (EverythingTarget) target() {}
func (Ambiguous) target()        {}
func (NoMatch) target()          {}

var (
	endsWithRoom = regexp.MustCompile(`\b(room|space)$`)
	roomSuffix   = regexp.MustCompile(` ?\b(room|space)$`)
)

// Match finds what a spoken name refers to. Exact names win over partial
// ones, and at the same strength a device wins over a space and a space
// over a mode, so "heater" picks the heater even in a room called "Heater
// room". Saying "room" or "space" after the name means the room. Two
// candidates of the same strength are Ambiguous. With modesOnly set, only
// modes are considered.
func Match(spoken string, home Home, modesOnly bool) Target {
	wanted := Normalize(spoken)
	if wanted == "" {
		return NoMatch{Spoken: spoken}
	}
	roomOnly := endsWithRoom.MatchString(wanted)
	roomName := roomSuffix.ReplaceAllString(wanted, "")

	exact := func(name string) bool { return Normalize(name) == wanted }
	partial := func(name string) bool {
		candidate := Normalize(name)
		return candidate != "" &&
			(strings.Contains(candidate, wanted) || strings.Contains(wanted, candidate))
	}
	spaceExact := func(name string) bool {
		return exact(name) || Normalize(name) == roomName
	}

	if modesOnly {
		for _, tier := range [][]Target{
			matchModes(home.Modes, wanted, true),
			matchModes(home.Modes, wanted, false),
		} {
			if t, ok := pick(spoken, tier); ok {
				return t
			}
		}
		return NoMatch{Spoken: spoken}
	}

	devices := func(fits func(string) bool) []Target {
		if roomOnly {
			return nil
		}
		var tier []Target
		for _, d := range home.Devices {
			if fits(d.Name) {
				tier = append(tier, DeviceTarget{Device: d})
			}
		}
		return tier
	}
	spaces := func(fits func(string) bool) []Target {
		var tier []Target
		for _, s := range home.Spaces {
			if fits(s.Name) {
				tier = append(tier, SpaceTarget{Space: s, Devices: home.DevicesIn(s)})
			}
		}
		return tier
	}
	for _, tier := range [][]Target{
		devices(exact),
		spaces(spaceExact),
		matchModes(home.Modes, wanted, true),
		devices(partial),
		spaces(partial),
		matchModes(home.Modes, wanted, false),
	} {
		if t, ok := pick(spoken, tier); ok {
			return t
		}
	}
	return NoMatch{Spoken: spoken}
}

// Resolve matches a parsed command to the home. A mode asked for in a space
// ("set the kitchen to Away mode") is looked for among that space's modes.
func Resolve(command Command, home Home) Target {
	switch command.Kind {
	case Everything:
		return EverythingTarget{Devices: home.Devices}
	case ModeKind:
		inSpace := home
		inSpace.Modes = modesInSpace(home.Modes, command.Space)
		return Match(command.Target, inSpace, true)
	default:
		return Match(command.Target, home, false)
	}
}

// matchModes returns the saved modes a name fits. An exact match on a voice
// name wins outright: when "Boss Time" exists as a general mode and in two
// spaces, only the general one is said "Boss Time".
func matchModes(modes []Mode, wanted string, exactMatch bool) []Target {
	var saved []Mode
	for _, m := range modes {
		if !m.Draft && m.ID != "" {
			saved = append(saved, m)
		}
	}
	if exactMatch {
		var byVoiceName []Target
		for _, m := range saved {
			if m.VoiceName != "" && Normalize(m.VoiceName) == wanted {
				byVoiceName = append(byVoiceName, ModeTarget{Mode: m})
			}
		}
		if len(byVoiceName) > 0 {
			return byVoiceName
		}
	}
	var tier []Target
	for _, m := range saved {
		for _, name := range []string{m.Name, m.VoiceName} {
			candidate := Normalize(name)
			if candidate == "" {
				continue
			}
			fits := candidate == wanted
			if !exactMatch {
				fits = strings.Contains(candidate, wanted) || strings.Contains(wanted, candidate)
			}
			if fits {
				tier = append(tier, ModeTarget{Mode: m})
				break
			}
		}
	}
	return tier
}

// modesInSpace keeps the modes of the named space: an exact space name
// first, else one that contains or is contained in what was said.
func modesInSpace(modes []Mode, space string) []Mode {
	if space == "" {
		return modes
	}
	var exact, partial []Mode
	for _, m := range modes {
		if m.SpaceName == "" {
			continue
		}
		name := Normalize(m.SpaceName)
		if name == space {
			exact = append(exact, m)
		}
		if strings.Contains(name, space) || strings.Contains(space, name) {
			partial = append(partial, m)
		}
	}
	if len(exact) > 0 {
		return exact
	}
	return partial
}

func pick(spoken string, tier []Target) (Target, bool) {
	switch len(tier) {
	case 0:
		return nil, false
	case 1:
		return tier[0], true
	default:
		return Ambiguous{Spoken: spoken, Options: tier}, true
	}
}
