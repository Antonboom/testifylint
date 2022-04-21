package main

import "testing"

func TestExpandCheck(t *testing.T) {
	cases := []struct {
		check           Check
		pkg             string
		dynamicArgValue string
		expected        string
	}{
		{
			check: Check{
				Fn:              "Len",
				Args:            []string{"t", "%s", "0"},
				DynamicArgIndex: 1,
				ReportedMsg:     "use %s.Empty",
			},
			pkg:             "assert",
			dynamicArgValue: "vv",
			expected: `assert.Len(t, vv, 0) // want "use assert.Empty"
assert.Len(t, vv, 0, "msg") // want "use assert.Empty"
assert.Len(t, vv, 0, "msg with arg %d", 42) // want "use assert.Empty"
assert.Lenf(t, vv, 0, "msg") // want "use assert.Emptyf"
assert.Lenf(t, vv, 0, "msg with arg %d", 42) // want "use assert.Emptyf"`,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			got := ExpandCheck(tt.check, tt.pkg, tt.dynamicArgValue)
			if got != tt.expected {
				t.Fatalf("\n%v\n!=\n%v", got, tt.expected)
			}
		})
	}
}
