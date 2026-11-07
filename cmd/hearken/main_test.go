package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMatchFromTheCommandLine(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"match", "--home", "../../examples/home.json", "set the kitchen to away mode"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Target struct {
			Mode struct{ ID string } `json:"mode"`
		} `json:"target"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Target.Mode.ID != "kitchen-away" {
		t.Fatalf("got %s", out.String())
	}
}

func TestNotACommandIsAnError(t *testing.T) {
	if err := run([]string{"parse", "what is the weather"}, &bytes.Buffer{}); err == nil {
		t.Fatal("want an error")
	}
}
