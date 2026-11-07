package hearken

import (
	"regexp"
	"strings"
)

// Action is what a command asks for.
type Action string

const (
	TurnOn  Action = "on"
	TurnOff Action = "off"
	Status  Action = "status"
)

// Kind says what the spoken target names.
type Kind string

const (
	// DeviceOrSpace is a device or a room, decided when the name is matched.
	DeviceOrSpace Kind = "device_or_space"
	// ModeKind is a scene such as "Home" or "Movie night".
	ModeKind Kind = "mode"
	// Everything is every device in the home.
	Everything Kind = "everything"
)

// Command is what was asked for, before the spoken name is matched to
// anything: "turn off the living room heater" is TurnOff with the target
// "living room heater".
type Command struct {
	Action Action `json:"action"`
	Kind   Kind   `json:"kind"`
	Target string `json:"target"`
	// Space is the room a mode was asked for in ("set the kitchen to Away
	// mode"), so only that room's modes are considered. Empty otherwise.
	Space string `json:"space,omitempty"`
}

// Parser reads spoken commands. The zero value has no wake phrase; use
// NewParser to strip one. A Parser is safe to share between goroutines.
type Parser struct {
	wake *regexp.Regexp
}

// NewParser returns a parser that ignores a wake phrase at the start of what
// was said. wake is a regular expression matched against normalised text
// (see Normalize), so it can also cover the ways a recogniser mishears the
// phrase, for example `^hey (home )?(flex|flicks) `. Empty means none.
func NewParser(wake string) (*Parser, error) {
	if wake == "" {
		return &Parser{}, nil
	}
	compiled, err := regexp.Compile(wake)
	if err != nil {
		return nil, err
	}
	return &Parser{wake: compiled}, nil
}

var (
	leadingHey      = regexp.MustCompile(`^hey `)
	leadingPolite   = regexp.MustCompile(`^(please |can you |could you )+`)
	trailingPlease  = regexp.MustCompile(` please$`)
	activateMode    = regexp.MustCompile(`^(?:activate|start|enable|run) (.+?) mode$`)
	statusQuery     = regexp.MustCompile(`^(?:is|are) (.+?) (?:switched |turned )?(?:on|off)$`)
	powerFirst      = regexp.MustCompile(`^(?:turn|switch|power) (on|off) (.+)$`)
	powerLast       = regexp.MustCompile(`^(?:turn|switch|power) (.+) (on|off)$`)
	modeInSpace     = regexp.MustCompile(`^(?:set|put|switch|change|turn) (?:the )?(.+?) (?:to|in|into|on) (?:the )?(.+?) mode$`)
	targetFiller    = regexp.MustCompile(`^(?:all |every )?(?:the )?(?:devices |things |lights )?(?:in |at )?(?:the )?`)
	everythingIn    = regexp.MustCompile(`^everything (?:in |at )?(?:the )?`)
	everythingNames = map[string]bool{
		"everything": true, "all": true, "all devices": true,
		"all the devices": true, "everything in the house": true,
		"every device": true,
	}
)

// Parse reads one spoken command. It reports false when what was said is
// not a command it understands.
func (p *Parser) Parse(spoken string) (Command, bool) {
	if command, ok := p.parseModeInSpace(spoken); ok {
		return command, true
	}
	text := p.stripWake(Normalize(spoken))
	text = leadingHey.ReplaceAllString(text, "")
	text = stripPolite(text)

	if m := activateMode.FindStringSubmatch(text); m != nil {
		return Command{Action: TurnOn, Kind: ModeKind, Target: m[1]}, true
	}
	if m := statusQuery.FindStringSubmatch(text); m != nil {
		return command(Status, m[1])
	}
	if m := powerFirst.FindStringSubmatch(text); m != nil {
		return command(power(m[1]), m[2])
	}
	if m := powerLast.FindStringSubmatch(text); m != nil {
		return command(power(m[2]), m[1])
	}
	return Command{}, false
}

// SpokenText is what was said in the form shortcuts compare: normalised,
// without the wake phrase or "please".
func (p *Parser) SpokenText(spoken string) string {
	return stripPolite(p.stripWake(Normalize(spoken)))
}

// parseModeInSpace reads "set Space A to Home mode" before number words are
// turned into digits, since that reads "to" as 2 and joins it to the room's
// number.
func (p *Parser) parseModeInSpace(spoken string) (Command, bool) {
	words := strings.ToLower(spoken)
	words = notWordOrSpace.ReplaceAllString(words, " ")
	words = strings.TrimSpace(runsOfSpace.ReplaceAllString(words, " "))
	words = stripPolite(p.stripWake(words))

	m := modeInSpace.FindStringSubmatch(words)
	if m == nil {
		return Command{}, false
	}
	space, mode := Normalize(m[1]), Normalize(m[2])
	if space == "" || mode == "" {
		return Command{}, false
	}
	return Command{Action: TurnOn, Kind: ModeKind, Target: mode, Space: space}, true
}

func (p *Parser) stripWake(text string) string {
	if p.wake == nil {
		return text
	}
	if loc := p.wake.FindStringIndex(text); loc != nil && loc[0] == 0 {
		return text[loc[1]:]
	}
	return text
}

func stripPolite(text string) string {
	text = leadingPolite.ReplaceAllString(text, "")
	return trailingPlease.ReplaceAllString(text, "")
}

func power(word string) Action {
	if word == "on" {
		return TurnOn
	}
	return TurnOff
}

func command(action Action, target string) (Command, bool) {
	trimmed := strings.TrimSpace(target)
	if everythingNames[trimmed] {
		return Command{Action: action, Kind: Everything, Target: trimmed}, true
	}
	name := targetFiller.ReplaceAllString(target, "")
	name = strings.TrimSpace(everythingIn.ReplaceAllString(name, ""))
	kind := DeviceOrSpace
	if strings.HasSuffix(name, " mode") {
		kind = ModeKind
		name = strings.TrimSpace(strings.TrimSuffix(name, " mode"))
	}
	if name == "" {
		return Command{}, false
	}
	return Command{Action: action, Kind: kind, Target: name}, true
}
