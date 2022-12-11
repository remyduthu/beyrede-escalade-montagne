package beyredeescalademontagne

import (
	"testing"
)

func TestStateValidate(t *testing.T) {
	cases := map[string]struct {
		state       string
		expectError bool
	}{
		"busy":    {state: string(busyState), expectError: false},
		"closed":  {state: string(closedState), expectError: false},
		"opened":  {state: string(openedState), expectError: false},
		"unknown": {state: "foo", expectError: true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := (state(c.state).validate() != nil); got != c.expectError {
				t.Fatalf("expected error: %t, got: %t", c.expectError, got)
			}
		})
	}
}
