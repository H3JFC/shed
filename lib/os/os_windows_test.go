//go:build windows

package os

import "testing"

func TestDetect_Windows(t *testing.T) {
	t.Parallel()

	if Detect() != Windows {
		t.Errorf("Detect() = %v, want %v", Detect(), ToOS(windows))
	}
}
