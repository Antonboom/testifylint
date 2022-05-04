package main

import "testing"

func TestCheckerExpander_Expand(t *testing.T) {
	cases := []struct {
		check          Check
		selector       string
		argValues      []any
		withoutTArg    bool
		withoutFFuncs  bool
		expected       string
		expectedGolden string
	}{
		// Basic.
		{
			check: Check{
				Fn:            "Len",
				Argsf:         "%s, 0",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "Empty",
				ProposedArgsf: "%s",
			},
			selector:  "assert",
			argValues: []any{"vv"},

			expected: `assert.Len(t, vv, 0) // want "use assert.Empty"
assert.Len(t, vv, 0, "msg") // want "use assert.Empty"
assert.Len(t, vv, 0, "msg with arg %d", 42) // want "use assert.Empty"
assert.Lenf(t, vv, 0, "msg") // want "use assert.Emptyf"
assert.Lenf(t, vv, 0, "msg with arg %d", 42) // want "use assert.Emptyf"`,

			expectedGolden: `assert.Empty(t, vv) // want "use assert.Empty"
assert.Empty(t, vv, "msg") // want "use assert.Empty"
assert.Empty(t, vv, "msg with arg %d", 42) // want "use assert.Empty"
assert.Emptyf(t, vv, "msg") // want "use assert.Emptyf"
assert.Emptyf(t, vv, "msg with arg %d", 42) // want "use assert.Emptyf"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "%s, %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:  "assert",
			argValues: []any{"42.42", "flNum"},

			expected: `assert.Equal(t, 42.42, flNum) // want "use assert.InDelta"
assert.Equal(t, 42.42, flNum, "msg") // want "use assert.InDelta"
assert.Equal(t, 42.42, flNum, "msg with arg %d", 42) // want "use assert.InDelta"
assert.Equalf(t, 42.42, flNum, "msg") // want "use assert.InDeltaf"
assert.Equalf(t, 42.42, flNum, "msg with arg %d", 42) // want "use assert.InDeltaf"`,

			expectedGolden: `assert.InDelta(t, 42.42, flNum, 0.0001) // want "use assert.InDelta"
assert.InDelta(t, 42.42, flNum, 0.0001, "msg") // want "use assert.InDelta"
assert.InDelta(t, 42.42, flNum, 0.0001, "msg with arg %d", 42) // want "use assert.InDelta"
assert.InDeltaf(t, 42.42, flNum, 0.0001, "msg") // want "use assert.InDeltaf"
assert.InDeltaf(t, 42.42, flNum, 0.0001, "msg with arg %d", 42) // want "use assert.InDeltaf"`,
		},
		{
			check: Check{
				Fn:            "True",
				Argsf:         "%s == %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:  "assert",
			argValues: []any{"a", "b"},

			expected: `assert.True(t, a == b) // want "use assert.InDelta"
assert.True(t, a == b, "msg") // want "use assert.InDelta"
assert.True(t, a == b, "msg with arg %d", 42) // want "use assert.InDelta"
assert.Truef(t, a == b, "msg") // want "use assert.InDeltaf"
assert.Truef(t, a == b, "msg with arg %d", 42) // want "use assert.InDeltaf"`,

			expectedGolden: `assert.InDelta(t, a, b, 0.0001) // want "use assert.InDelta"
assert.InDelta(t, a, b, 0.0001, "msg") // want "use assert.InDelta"
assert.InDelta(t, a, b, 0.0001, "msg with arg %d", 42) // want "use assert.InDelta"
assert.InDeltaf(t, a, b, 0.0001, "msg") // want "use assert.InDeltaf"
assert.InDeltaf(t, a, b, 0.0001, "msg with arg %d", 42) // want "use assert.InDeltaf"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "result, %s",
				ReportMsgf:    "need to reverse actual and expected values",
				ProposedArgsf: "%s, result",
			},
			selector:  "assert",
			argValues: []any{`"expected string"`},

			expected: `assert.Equal(t, result, "expected string") // want "need to reverse actual and expected values"
assert.Equal(t, result, "expected string", "msg") // want "need to reverse actual and expected values"
assert.Equal(t, result, "expected string", "msg with arg %d", 42) // want "need to reverse actual and expected values"
assert.Equalf(t, result, "expected string", "msg") // want "need to reverse actual and expected values"
assert.Equalf(t, result, "expected string", "msg with arg %d", 42) // want "need to reverse actual and expected values"`,

			expectedGolden: `assert.Equal(t, "expected string", result) // want "need to reverse actual and expected values"
assert.Equal(t, "expected string", result, "msg") // want "need to reverse actual and expected values"
assert.Equal(t, "expected string", result, "msg with arg %d", 42) // want "need to reverse actual and expected values"
assert.Equalf(t, "expected string", result, "msg") // want "need to reverse actual and expected values"
assert.Equalf(t, "expected string", result, "msg with arg %d", 42) // want "need to reverse actual and expected values"`,
		},
		{
			check: Check{
				Fn:    "InDelta",
				Argsf: "%s, %s, 0.0001",
			},
			selector:  "assert",
			argValues: []any{"s.c", "h.Calculate()"},

			expected: `assert.InDelta(t, s.c, h.Calculate(), 0.0001)
assert.InDelta(t, s.c, h.Calculate(), 0.0001, "msg")
assert.InDelta(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)
assert.InDeltaf(t, s.c, h.Calculate(), 0.0001, "msg")
assert.InDeltaf(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,

			expectedGolden: `assert.InDelta(t, s.c, h.Calculate(), 0.0001)
assert.InDelta(t, s.c, h.Calculate(), 0.0001, "msg")
assert.InDelta(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)
assert.InDeltaf(t, s.c, h.Calculate(), 0.0001, "msg")
assert.InDeltaf(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,
		},

		// Without f-funcs.
		{
			check: Check{
				Fn:            "Len",
				Argsf:         "%s, 0",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "Empty",
				ProposedArgsf: "%s",
			},
			selector:      "require",
			argValues:     []any{"vv"},
			withoutFFuncs: true,

			expected: `require.Len(t, vv, 0) // want "use require.Empty"
require.Len(t, vv, 0, "msg") // want "use require.Empty"
require.Len(t, vv, 0, "msg with arg %d", 42) // want "use require.Empty"`,

			expectedGolden: `require.Empty(t, vv) // want "use require.Empty"
require.Empty(t, vv, "msg") // want "use require.Empty"
require.Empty(t, vv, "msg with arg %d", 42) // want "use require.Empty"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "%s, %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:      "require",
			argValues:     []any{"42.42", "flNum"},
			withoutFFuncs: true,

			expected: `require.Equal(t, 42.42, flNum) // want "use require.InDelta"
require.Equal(t, 42.42, flNum, "msg") // want "use require.InDelta"
require.Equal(t, 42.42, flNum, "msg with arg %d", 42) // want "use require.InDelta"`,

			expectedGolden: `require.InDelta(t, 42.42, flNum, 0.0001) // want "use require.InDelta"
require.InDelta(t, 42.42, flNum, 0.0001, "msg") // want "use require.InDelta"
require.InDelta(t, 42.42, flNum, 0.0001, "msg with arg %d", 42) // want "use require.InDelta"`,
		},
		{
			check: Check{
				Fn:            "True",
				Argsf:         "%s == %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:      "require",
			argValues:     []any{"a", "b"},
			withoutFFuncs: true,

			expected: `require.True(t, a == b) // want "use require.InDelta"
require.True(t, a == b, "msg") // want "use require.InDelta"
require.True(t, a == b, "msg with arg %d", 42) // want "use require.InDelta"`,

			expectedGolden: `require.InDelta(t, a, b, 0.0001) // want "use require.InDelta"
require.InDelta(t, a, b, 0.0001, "msg") // want "use require.InDelta"
require.InDelta(t, a, b, 0.0001, "msg with arg %d", 42) // want "use require.InDelta"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "result, %s",
				ReportMsgf:    "need to reverse actual and expected values",
				ProposedArgsf: "%s, result",
			},
			selector:      "require",
			argValues:     []any{`"expected string"`},
			withoutFFuncs: true,

			expected: `require.Equal(t, result, "expected string") // want "need to reverse actual and expected values"
require.Equal(t, result, "expected string", "msg") // want "need to reverse actual and expected values"
require.Equal(t, result, "expected string", "msg with arg %d", 42) // want "need to reverse actual and expected values"`,

			expectedGolden: `require.Equal(t, "expected string", result) // want "need to reverse actual and expected values"
require.Equal(t, "expected string", result, "msg") // want "need to reverse actual and expected values"
require.Equal(t, "expected string", result, "msg with arg %d", 42) // want "need to reverse actual and expected values"`,
		},
		{
			check: Check{
				Fn:            "InDelta",
				Argsf:         "%s, %s, 0.0001",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:      "require",
			argValues:     []any{"s.c", "h.Calculate()"},
			withoutFFuncs: true,

			expected: `require.InDelta(t, s.c, h.Calculate(), 0.0001)
require.InDelta(t, s.c, h.Calculate(), 0.0001, "msg")
require.InDelta(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,

			expectedGolden: `require.InDelta(t, s.c, h.Calculate(), 0.0001)
require.InDelta(t, s.c, h.Calculate(), 0.0001, "msg")
require.InDelta(t, s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,
		},

		// Without t arg.
		{
			check: Check{
				Fn:            "Len",
				Argsf:         "%s, 0",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "Empty",
				ProposedArgsf: "%s",
			},
			selector:    "assertObj",
			argValues:   []any{"vv"},
			withoutTArg: true,

			expected: `assertObj.Len(vv, 0) // want "use assertObj.Empty"
assertObj.Len(vv, 0, "msg") // want "use assertObj.Empty"
assertObj.Len(vv, 0, "msg with arg %d", 42) // want "use assertObj.Empty"
assertObj.Lenf(vv, 0, "msg") // want "use assertObj.Emptyf"
assertObj.Lenf(vv, 0, "msg with arg %d", 42) // want "use assertObj.Emptyf"`,

			expectedGolden: `assertObj.Empty(vv) // want "use assertObj.Empty"
assertObj.Empty(vv, "msg") // want "use assertObj.Empty"
assertObj.Empty(vv, "msg with arg %d", 42) // want "use assertObj.Empty"
assertObj.Emptyf(vv, "msg") // want "use assertObj.Emptyf"
assertObj.Emptyf(vv, "msg with arg %d", 42) // want "use assertObj.Emptyf"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "%s, %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:    "assertObj",
			argValues:   []any{"42.42", "flNum"},
			withoutTArg: true,

			expected: `assertObj.Equal(42.42, flNum) // want "use assertObj.InDelta"
assertObj.Equal(42.42, flNum, "msg") // want "use assertObj.InDelta"
assertObj.Equal(42.42, flNum, "msg with arg %d", 42) // want "use assertObj.InDelta"
assertObj.Equalf(42.42, flNum, "msg") // want "use assertObj.InDeltaf"
assertObj.Equalf(42.42, flNum, "msg with arg %d", 42) // want "use assertObj.InDeltaf"`,

			expectedGolden: `assertObj.InDelta(42.42, flNum, 0.0001) // want "use assertObj.InDelta"
assertObj.InDelta(42.42, flNum, 0.0001, "msg") // want "use assertObj.InDelta"
assertObj.InDelta(42.42, flNum, 0.0001, "msg with arg %d", 42) // want "use assertObj.InDelta"
assertObj.InDeltaf(42.42, flNum, 0.0001, "msg") // want "use assertObj.InDeltaf"
assertObj.InDeltaf(42.42, flNum, 0.0001, "msg with arg %d", 42) // want "use assertObj.InDeltaf"`,
		},
		{
			check: Check{
				Fn:            "True",
				Argsf:         "%s == %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:    "assertObj",
			argValues:   []any{"d", "c"},
			withoutTArg: true,

			expected: `assertObj.True(d == c) // want "use assertObj.InDelta"
assertObj.True(d == c, "msg") // want "use assertObj.InDelta"
assertObj.True(d == c, "msg with arg %d", 42) // want "use assertObj.InDelta"
assertObj.Truef(d == c, "msg") // want "use assertObj.InDeltaf"
assertObj.Truef(d == c, "msg with arg %d", 42) // want "use assertObj.InDeltaf"`,

			expectedGolden: `assertObj.InDelta(d, c, 0.0001) // want "use assertObj.InDelta"
assertObj.InDelta(d, c, 0.0001, "msg") // want "use assertObj.InDelta"
assertObj.InDelta(d, c, 0.0001, "msg with arg %d", 42) // want "use assertObj.InDelta"
assertObj.InDeltaf(d, c, 0.0001, "msg") // want "use assertObj.InDeltaf"
assertObj.InDeltaf(d, c, 0.0001, "msg with arg %d", 42) // want "use assertObj.InDeltaf"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "res, %s",
				ReportMsgf:    "need to reverse actual and expected values",
				ProposedArgsf: "%s, res",
			},
			selector:    "assertObj",
			argValues:   []any{`"expected string"`},
			withoutTArg: true,

			expected: `assertObj.Equal(res, "expected string") // want "need to reverse actual and expected values"
assertObj.Equal(res, "expected string", "msg") // want "need to reverse actual and expected values"
assertObj.Equal(res, "expected string", "msg with arg %d", 42) // want "need to reverse actual and expected values"
assertObj.Equalf(res, "expected string", "msg") // want "need to reverse actual and expected values"
assertObj.Equalf(res, "expected string", "msg with arg %d", 42) // want "need to reverse actual and expected values"`,

			expectedGolden: `assertObj.Equal("expected string", res) // want "need to reverse actual and expected values"
assertObj.Equal("expected string", res, "msg") // want "need to reverse actual and expected values"
assertObj.Equal("expected string", res, "msg with arg %d", 42) // want "need to reverse actual and expected values"
assertObj.Equalf("expected string", res, "msg") // want "need to reverse actual and expected values"
assertObj.Equalf("expected string", res, "msg with arg %d", 42) // want "need to reverse actual and expected values"`,
		},
		{
			check: Check{
				Fn:    "InDelta",
				Argsf: "%s, %s, 0.0001",
			},
			selector:    "assertObj",
			argValues:   []any{"s.c", "h.Calculate()"},
			withoutTArg: true,

			expected: `assertObj.InDelta(s.c, h.Calculate(), 0.0001)
assertObj.InDelta(s.c, h.Calculate(), 0.0001, "msg")
assertObj.InDelta(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)
assertObj.InDeltaf(s.c, h.Calculate(), 0.0001, "msg")
assertObj.InDeltaf(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,

			expectedGolden: `assertObj.InDelta(s.c, h.Calculate(), 0.0001)
assertObj.InDelta(s.c, h.Calculate(), 0.0001, "msg")
assertObj.InDelta(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)
assertObj.InDeltaf(s.c, h.Calculate(), 0.0001, "msg")
assertObj.InDeltaf(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,
		},

		// Without f-funcs and t arg.
		{
			check: Check{
				Fn:            "Len",
				Argsf:         "%s, 0",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "Empty",
				ProposedArgsf: "%s",
			},
			selector:      "s.Require()",
			argValues:     []any{"users"},
			withoutFFuncs: true,
			withoutTArg:   true,

			expected: `s.Require().Len(users, 0) // want "use s.Require().Empty"
s.Require().Len(users, 0, "msg") // want "use s.Require().Empty"
s.Require().Len(users, 0, "msg with arg %d", 42) // want "use s.Require().Empty"`,

			expectedGolden: `s.Require().Empty(users) // want "use s.Require().Empty"
s.Require().Empty(users, "msg") // want "use s.Require().Empty"
s.Require().Empty(users, "msg with arg %d", 42) // want "use s.Require().Empty"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "%s, %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:      "s.Require()",
			argValues:     []any{"floatOp()", "pi"},
			withoutFFuncs: true,
			withoutTArg:   true,

			expected: `s.Require().Equal(floatOp(), pi) // want "use s.Require().InDelta"
s.Require().Equal(floatOp(), pi, "msg") // want "use s.Require().InDelta"
s.Require().Equal(floatOp(), pi, "msg with arg %d", 42) // want "use s.Require().InDelta"`,

			expectedGolden: `s.Require().InDelta(floatOp(), pi, 0.0001) // want "use s.Require().InDelta"
s.Require().InDelta(floatOp(), pi, 0.0001, "msg") // want "use s.Require().InDelta"
s.Require().InDelta(floatOp(), pi, 0.0001, "msg with arg %d", 42) // want "use s.Require().InDelta"`,
		},
		{
			check: Check{
				Fn:            "True",
				Argsf:         "%s == %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InDelta",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			selector:      "s.Require()",
			argValues:     []any{"a", "b"},
			withoutFFuncs: true,
			withoutTArg:   true,

			expected: `s.Require().True(a == b) // want "use s.Require().InDelta"
s.Require().True(a == b, "msg") // want "use s.Require().InDelta"
s.Require().True(a == b, "msg with arg %d", 42) // want "use s.Require().InDelta"`,

			expectedGolden: `s.Require().InDelta(a, b, 0.0001) // want "use s.Require().InDelta"
s.Require().InDelta(a, b, 0.0001, "msg") // want "use s.Require().InDelta"
s.Require().InDelta(a, b, 0.0001, "msg with arg %d", 42) // want "use s.Require().InDelta"`,
		},
		{
			check: Check{
				Fn:            "Equal",
				Argsf:         "count, %s",
				ReportMsgf:    "need to reverse actual and expected values",
				ProposedArgsf: "%s, count",
			},
			selector:      "s.Require()",
			argValues:     []any{`100`},
			withoutFFuncs: true,
			withoutTArg:   true,

			expected: `s.Require().Equal(count, 100) // want "need to reverse actual and expected values"
s.Require().Equal(count, 100, "msg") // want "need to reverse actual and expected values"
s.Require().Equal(count, 100, "msg with arg %d", 42) // want "need to reverse actual and expected values"`,

			expectedGolden: `s.Require().Equal(100, count) // want "need to reverse actual and expected values"
s.Require().Equal(100, count, "msg") // want "need to reverse actual and expected values"
s.Require().Equal(100, count, "msg with arg %d", 42) // want "need to reverse actual and expected values"`,
		},
		{
			check: Check{
				Fn:    "InDelta",
				Argsf: "%s, %s, 0.0001",
			},
			selector:      "s.Require()",
			argValues:     []any{"s.c", "h.Calculate()"},
			withoutFFuncs: true,
			withoutTArg:   true,

			expected: `s.Require().InDelta(s.c, h.Calculate(), 0.0001)
s.Require().InDelta(s.c, h.Calculate(), 0.0001, "msg")
s.Require().InDelta(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,

			expectedGolden: `s.Require().InDelta(s.c, h.Calculate(), 0.0001)
s.Require().InDelta(s.c, h.Calculate(), 0.0001, "msg")
s.Require().InDelta(s.c, h.Calculate(), 0.0001, "msg with arg %d", 42)`,
		},

		// Special cases.
		{
			check: Check{
				Fn:         "Error",
				Argsf:      "err, errSentinel",
				ReportMsgf: "invalid usage of %[1]s.Error, use %[1]s.%[2]s instead",
				ProposedFn: "ErrorIs",
			},
			selector:      "r",
			argValues:     nil,
			withoutFFuncs: true,

			expected: `r.Error(t, err, errSentinel) // want "invalid usage of r.Error, use r.ErrorIs instead"
r.Error(t, err, errSentinel, "msg") // want "invalid usage of r.Error, use r.ErrorIs instead"
r.Error(t, err, errSentinel, "msg with arg %d", 42) // want "invalid usage of r.Error, use r.ErrorIs instead"`,

			expectedGolden: `r.ErrorIs(t, err, errSentinel) // want "invalid usage of r.Error, use r.ErrorIs instead"
r.ErrorIs(t, err, errSentinel, "msg") // want "invalid usage of r.Error, use r.ErrorIs instead"
r.ErrorIs(t, err, errSentinel, "msg with arg %d", 42) // want "invalid usage of r.Error, use r.ErrorIs instead"`,
		},
	}

	for _, tt := range cases {
		newExpanderFromCase := func() *CheckerExpander {
			expander := NewCheckerExpander()

			if tt.withoutTArg {
				expander = expander.WithoutTArg()
			}
			if tt.withoutFFuncs {
				expander = expander.WithoutFFuncs()
			}

			return expander
		}

		t.Run("", func(t *testing.T) {
			t.Run("errored", func(t *testing.T) {
				got := newExpanderFromCase().Expand(tt.check, tt.selector, tt.argValues)
				assertEqualStrings(t, tt.expected, got)
			})

			t.Run("golden", func(t *testing.T) {
				got := newExpanderFromCase().AsGolden().Expand(tt.check, tt.selector, tt.argValues)
				assertEqualStrings(t, tt.expectedGolden, got)
			})
		})
	}
}

func assertEqualStrings(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Fatalf("\nactual\n%v\n!= expected\n%v", actual, expected)
	}
}
