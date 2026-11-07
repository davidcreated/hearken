package hearken

import "testing"

func TestShortcutSpec(t *testing.T) {
	var spec struct {
		Wake          string     `json:"wake"`
		Shortcuts     []Shortcut `json:"shortcuts"`
		Nicknames     []Nickname `json:"nicknames"`
		ShortcutCases []struct {
			Say    string  `json:"say"`
			Expect *string `json:"expect"`
		} `json:"shortcutCases"`
		NicknameCases []struct {
			Say    string  `json:"say"`
			Expect *string `json:"expect"`
		} `json:"nicknameCases"`
	}
	readSpec(t, "shortcuts.json", &spec)
	parser, err := NewParser(spec.Wake)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(spec.Shortcuts[0].Actions); got != 3 {
		t.Fatalf("every action of a shortcut is read, got %d", got)
	}
	for _, c := range spec.ShortcutCases {
		got, ok := parser.FindShortcut(c.Say, spec.Shortcuts)
		checkFound(t, c.Say, got.ID, ok, c.Expect)
	}
	for _, c := range spec.NicknameCases {
		got, ok := FindNickname(c.Say, spec.Nicknames)
		checkFound(t, c.Say, got.Alias, ok, c.Expect)
	}
}

func checkFound(t *testing.T, said, got string, found bool, want *string) {
	t.Helper()
	switch {
	case want == nil && found:
		t.Errorf("%q: found %q, want nothing", said, got)
	case want != nil && !found:
		t.Errorf("%q: found nothing, want %q", said, *want)
	case want != nil && got != *want:
		t.Errorf("%q: found %q, want %q", said, got, *want)
	}
}
