//go:build darwin

package os

import "testing"

func TestDetect_Darwin(t *testing.T) {
	t.Parallel()

	if Detect() != Darwin {
		t.Errorf("Detect() = %v, want %v", Detect(), ToOS(darwin))
	}
}
