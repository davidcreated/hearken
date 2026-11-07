package hearken

import "strings"

// ShortcutActionKind is what one step of a shortcut does.
type ShortcutActionKind string

const (
	SwitchDevice ShortcutActionKind = "switch_device"
	SwitchSpace  ShortcutActionKind = "switch_space"
	TurnOnMode   ShortcutActionKind = "turn_on_mode"
)

// ShortcutAction is one step of a shortcut. It keeps the ID it acts on and
// the name shown to the user, so renaming a device doesn't break it.
type ShortcutAction struct {
	Kind ShortcutActionKind `json:"kind"`
	ID   string             `json:"id"`
	Name string             `json:"name"`
	// On is the state to switch to; unused for TurnOnMode.
	On bool `json:"on,omitempty"`
}

// Shortcut is a phrase of the user's own, such as "movie time", and the
// steps it runs.
type Shortcut struct {
	ID      string           `json:"id"`
	Phrase  string           `json:"phrase"`
	Actions []ShortcutAction `json:"actions"`
}

// FindShortcut returns the shortcut whose phrase is exactly what was said,
// ignoring case, punctuation, the wake phrase and "please".
func (p *Parser) FindShortcut(spoken string, shortcuts []Shortcut) (Shortcut, bool) {
	said := p.SpokenText(spoken)
	for _, s := range shortcuts {
		if p.SpokenText(s.Phrase) == said {
			return s, true
		}
	}
	return Shortcut{}, false
}

// NicknameTarget is what a nickname stands for.
type NicknameTarget string

const (
	NicknameDevice NicknameTarget = "device"
	NicknameSpace  NicknameTarget = "space"
)

// Nickname is another name for a device or room, such as "AC" for the
// bedroom air conditioner.
type Nickname struct {
	Alias      string         `json:"alias"`
	Target     NicknameTarget `json:"target"`
	TargetID   string         `json:"targetId"`
	TargetName string         `json:"targetName"`
}

// FindNickname returns the nickname that was said. Nicknames are often
// initials, which a recogniser may hear as "A.C." or "a c", so they are
// compared with spaces removed.
func FindNickname(spoken string, nicknames []Nickname) (Nickname, bool) {
	said := aliasKey(spoken)
	for _, n := range nicknames {
		if aliasKey(n.Alias) == said {
			return n, true
		}
	}
	return Nickname{}, false
}

func aliasKey(alias string) string {
	return strings.ReplaceAll(Normalize(alias), " ", "")
}
