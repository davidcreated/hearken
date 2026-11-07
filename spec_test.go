package hearken

import (
	"encoding/json"
	"os"
	"testing"
)

func readSpec(t *testing.T, name string, into any) {
	t.Helper()
	raw, err := os.ReadFile("spec/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("spec/%s: %v", name, err)
	}
}
