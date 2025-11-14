//go:build linux

package os

import "testing"

func TestDetect_Linux(t *testing.T) {
	t.Parallel()

	if Detect() != Linux {
		t.Errorf("Detect() = %v, want %v", Detect(), ToOS(linux))
	}
}
