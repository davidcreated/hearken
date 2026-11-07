package hearken

import "testing"

func TestParseSpec(t *testing.T) {
	var spec struct {
		Wake  string `json:"wake"`
		Cases []struct {
			Say    string   `json:"say"`
			Expect *Command `json:"expect"`
		} `json:"cases"`
	}
	readSpec(t, "parse.json", &spec)
	parser, err := NewParser(spec.Wake)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range spec.Cases {
		got, ok := parser.Parse(c.Say)
		switch {
		case c.Expect == nil && ok:
			t.Errorf("%q: got %+v, want no command", c.Say, got)
		case c.Expect != nil && !ok:
			t.Errorf("%q: got no command, want %+v", c.Say, *c.Expect)
		case c.Expect != nil && got != *c.Expect:
			t.Errorf("%q: got %+v, want %+v", c.Say, got, *c.Expect)
		}
	}
}

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{
		"Test Twenty Two!":   "test 22",
		"Kitchen Home 2":     "kitchen home 2",
		"Lamp3":              "lamp 3",
		"  A.C.  ":           "a c",
		"one two":            "12",
		"forty five degrees": "45 degrees",
		"Café  Lights":       "café lights",
		"":                   "",
	} {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestZeroParserHasNoWakePhrase(t *testing.T) {
	var parser Parser
	got, ok := parser.Parse("turn on the kettle")
	if !ok || got.Target != "kettle" {
		t.Fatalf("got %+v, %v", got, ok)
	}
	if _, ok := parser.Parse("hey, turn on the kettle"); !ok {
		t.Fatal("a bare leading hey is dropped even without a wake pattern")
	}
	if got, ok := parser.Parse("hey ohmflex turn on the kettle"); ok {
		t.Fatalf("an unknown wake phrase is not stripped, got %+v", got)
	}
}

func TestNewParserRejectsABadPattern(t *testing.T) {
	if _, err := NewParser("("); err == nil {
		t.Fatal("want an error for an invalid pattern")
	}
}
