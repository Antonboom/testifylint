package main

import "testing"

func TestExpandCheck(t *testing.T) {
	cases := []struct {
		check     Check
		pkg       string
		argValues []any
		expected  string
	}{
		{
			check: Check{
				Fn:          "Len",
				ArgsTmpl:    "t, %s, 0",
				ReportedMsg: "use %s.Empty",
			},
			pkg:       "assert",
			argValues: []any{"vv"},
			expected: `assert.Len(t, vv, 0) // want "use assert.Empty"
assert.Len(t, vv, 0, "msg") // want "use assert.Empty"
assert.Len(t, vv, 0, "msg with arg %d", 42) // want "use assert.Empty"
assert.Lenf(t, vv, 0, "msg") // want "use assert.Emptyf"
assert.Lenf(t, vv, 0, "msg with arg %d", 42) // want "use assert.Emptyf"`,
		},
		{
			check: Check{
				Fn:          "Equal",
				ArgsTmpl:    "t, %s, %s",
				ReportedMsg: "use %s.InDelta",
			},
			pkg:       "require",
			argValues: []any{"42.42", "flNum"},
			expected: `require.Equal(t, 42.42, flNum) // want "use require.InDelta"
require.Equal(t, 42.42, flNum, "msg") // want "use require.InDelta"
require.Equal(t, 42.42, flNum, "msg with arg %d", 42) // want "use require.InDelta"
require.Equalf(t, 42.42, flNum, "msg") // want "use require.InDeltaf"
require.Equalf(t, 42.42, flNum, "msg with arg %d", 42) // want "use require.InDeltaf"`,
		},
		{
			check: Check{
				Fn:          "True",
				ArgsTmpl:    "t, %s == %s",
				ReportedMsg: "use %s.InDelta",
			},
			pkg:       "assert",
			argValues: []any{"a", "b"},
			expected: `assert.True(t, a == b) // want "use assert.InDelta"
assert.True(t, a == b, "msg") // want "use assert.InDelta"
assert.True(t, a == b, "msg with arg %d", 42) // want "use assert.InDelta"
assert.Truef(t, a == b, "msg") // want "use assert.InDeltaf"
assert.Truef(t, a == b, "msg with arg %d", 42) // want "use assert.InDeltaf"`,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			got := ExpandCheck(tt.check, tt.pkg, tt.argValues)
			if got != tt.expected {
				t.Fatalf("\n%v\n!=\n%v", got, tt.expected)
			}
		})
	}
}
