// terminfo package provides simple parser for compiled terminfo fields
// and additionally it provides trivial API for applying of capability
// of a terminal.
package terminfo

import (
	"os"
	"testing"
)

func TestParseTerminfo(t *testing.T) {
	os.Setenv("TERM", "gnome-256color")
	result, err := ParseTerminfoFromEnv()
	if err != nil {
		t.Error("ParseTerminfoFromEnv failed")
	}

	if result.Name != "gnome-256color" {
		t.Error("ParseTerminfoFromEnv failed (Name)")
	}

	if result.Bools["msgr"] != true {
		t.Error("ParseTerminfoFromEnv failed (msgr)")
	}

	if result.Bools["bce"] != true {
		t.Error("ParseTerminfoFromEnv failed (bce)")
	}

	if result.Bools["ccc"] != true {
		t.Error("ParseTerminfoFromEnv failed (ccc)")
	}

	if result.Numbers["cols"] != 80 {
		t.Error("ParseTerminfoFromEnv failed (cols)")
	}

	if result.Numbers["lines"] != 24 {
		t.Error("ParseTerminfoFromEnv failed (lines)")
	}

	capability, err := result.ApplyCapability("setab", 5)
	expected := []byte{27, 91, 52, 53, 109}

	if err != nil {
		t.Error("ApplyCapability failed")
	}

	for i, v := range []byte(capability) {
		if v != expected[i] {
			t.Error("ApplyCapability (setab) failed")
		}
	}
}
