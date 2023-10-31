package main

import "testing"

func TestNewAssertionExpander(t *testing.T) {
	cases := []struct {
		name           string
		selector       string
		tParam         string
		args           []any
		assrn          Assertion
		expected       string
		expectedGolden string
	}{
		{
			name:           "no report without t param",
			selector:       "s.Require()",
			assrn:          Assertion{Fn: "Equal", Argsf: "42.42, result"},
			expected:       `s.Require().Equal(42.42, result)`,
			expectedGolden: `s.Require().Equal(42.42, result)`,
		},
		{
			name:           "no report with t param",
			selector:       "assert",
			tParam:         "t",
			assrn:          Assertion{Fn: "Equal", Argsf: "42.42, result"},
			expected:       `assert.Equal(t, 42.42, result)`,
			expectedGolden: `assert.Equal(t, 42.42, result)`,
		},
		{
			name:           "no report with suite T param",
			selector:       "assert",
			tParam:         "s.T()",
			assrn:          Assertion{Fn: "Equal", Argsf: "42.42, result"},
			expected:       `assert.Equal(s.T(), 42.42, result)`,
			expectedGolden: `assert.Equal(s.T(), 42.42, result)`,
		},
		{
			name:           "no report args formatting",
			selector:       "s",
			assrn:          Assertion{Fn: "Equal", Argsf: "%s, %s"},
			args:           []any{"expected", "actual"},
			expected:       `s.Equal(expected, actual)`,
			expectedGolden: `s.Equal(expected, actual)`,
		},
		{
			name:     "report without formatting  proposed args only 1",
			selector: "s",
			args:     []any{"predicate"},
			assrn: Assertion{
				Fn:            "True",
				Argsf:         "%s == true",
				ReportMsgf:    "need to simplify the assertion",
				ProposedArgsf: "%s",
			},
			expected:       `s.True(predicate == true) // want "need to simplify the assertion"`,
			expectedGolden: `s.True(predicate) // want "need to simplify the assertion"`,
		},
		{
			name:     "report without formatting  proposed args only 2",
			selector: "assert",
			args:     []any{`"expected string"`},
			assrn: Assertion{
				Fn:            "Equal",
				Argsf:         "result, %s",
				ReportMsgf:    "need to reverse actual and expected values",
				ProposedArgsf: "%s, result",
			},
			expected:       `assert.Equal(result, "expected string") // want "need to reverse actual and expected values"`,
			expectedGolden: `assert.Equal("expected string", result) // want "need to reverse actual and expected values"`,
		},
		{
			name:     "report without formatting  proposed selector only",
			selector: "s.Assert()",
			assrn: Assertion{
				Fn:               "Equal",
				Argsf:            "a, b",
				ReportMsgf:       "need to simplify the assertion%.s%.s",
				ProposedSelector: "s",
			},
			expected:       `s.Assert().Equal(a, b) // want "need to simplify the assertion"`,
			expectedGolden: `s.Equal(a, b) // want "need to simplify the assertion"`,
		},
		{
			name:     "report without formatting and without proposals",
			selector: "assert",
			tParam:   "s.T()",
			assrn: Assertion{
				Fn:         "True",
				Argsf:      "b",
				ReportMsgf: "use suite API instead of package", // Don't need `%.s%.s, because no proposals.
			},
			expected:       `assert.True(s.T(), b) // want "use suite API instead of package"`,
			expectedGolden: `assert.True(s.T(), b) // want "use suite API instead of package"`,
		},
		{
			name:     "report with formatting  proposed selector only",
			selector: "s",
			assrn: Assertion{
				Fn:               "NoError",
				Argsf:            "err",
				ReportMsgf:       "use %s.NoError%.s",
				ProposedSelector: "s.Require()",
			},
			expected:       `s.NoError(err) // want "use s\\.Require\\(\\)\\.NoError"`,
			expectedGolden: `s.Require().NoError(err) // want "use s\\.Require\\(\\)\\.NoError"`,
		},
		{
			name:     "report with formatting  proposed fn and args",
			selector: "require",
			tParam:   "t",
			args:     []any{"42.42", "result"},
			assrn: Assertion{
				Fn:            "Equal",
				Argsf:         "%s, %s",
				ReportMsgf:    "use %s.%s",
				ProposedFn:    "InEpsilon",
				ProposedArgsf: "%s, %s, 0.0001",
			},
			expected:       `require.Equal(t, 42.42, result) // want "use require\\.InEpsilon"`,
			expectedGolden: `require.InEpsilon(t, 42.42, result, 0.0001) // want "use require\\.InEpsilon"`,
		},
		{
			name:     "report with positional formatting",
			selector: "require",
			tParam:   "t",
			assrn: Assertion{
				Fn:         "Error",
				Argsf:      "err, errSentinel",
				ReportMsgf: "invalid usage of %[1]s.Error, use %[1]s.%[2]s instead",
				ProposedFn: "ErrorIs",
			},
			expected:       `require.Error(t, err, errSentinel) // want "invalid usage of require\\.Error, use require\\.ErrorIs instead"`,
			expectedGolden: `require.ErrorIs(t, err, errSentinel) // want "invalid usage of require\\.Error, use require\\.ErrorIs instead"`,
		},

		// Assertion.WithoutReport examples.
		{
			name:     "assertion with report",
			selector: "suiteObj",
			args:     []any{"arr"},
			assrn: Assertion{
				Fn:               "Equal",
				Argsf:            "3, len(%s)",
				ReportMsgf:       "use %s.%s",
				ProposedSelector: "suiteObj.Require()",
				ProposedArgsf:    "%s, 3",
				ProposedFn:       "Len",
			},
			expected:       `suiteObj.Equal(3, len(arr)) // want "use suiteObj\\.Require\\(\\)\\.Len"`,
			expectedGolden: `suiteObj.Require().Len(arr, 3) // want "use suiteObj\\.Require\\(\\)\\.Len"`,
		},
		{
			name:     "assertion without report",
			selector: "suiteObj",
			args:     []any{"arr"},
			assrn: Assertion{
				Fn:               "Equal",
				Argsf:            "3, len(%s)",
				ReportMsgf:       "use %s.%s",
				ProposedSelector: "suiteObj.Require()",
				ProposedArgsf:    "%s, 3",
				ProposedFn:       "Len",
			}.WithoutReport(),
			expected:       `suiteObj.Equal(3, len(arr))`,
			expectedGolden: `suiteObj.Equal(3, len(arr))`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("errored", func(t *testing.T) {
				expander := NewAssertionExpander().NotFmtSingleMode()
				got := expander.Expand(tt.assrn, tt.selector, tt.tParam, tt.args)
				assertEqualStrings(t, tt.expected, got)
			})

			t.Run("golden", func(t *testing.T) {
				expander := NewAssertionExpander().NotFmtSingleMode().AsGolden()
				got := expander.Expand(tt.assrn, tt.selector, tt.tParam, tt.args)
				assertEqualStrings(t, tt.expectedGolden, got)
			})
		})
	}
}

func TestAssertionExpander_DifferentModes(t *testing.T) {
	const selector = "assert"

	assrn := Assertion{
		Fn:            "Len",
		Argsf:         "arr, 0",
		ReportMsgf:    "use %s.%s",
		ProposedFn:    "Empty",
		ProposedArgsf: "arr",
	}

	type modeMethod func(expander *AssertionExpander) *AssertionExpander

	cases := []struct {
		name           string
		mode           modeMethod
		expected       string
		expectedGolden string
	}{
		{
			name: "default extreme mode",
			mode: nil,
			expected: `assert.Len(t, arr, 0) // want "use assert\\.Empty"
assert.Lenf(t, arr, 0, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
			expectedGolden: `assert.Empty(t, arr) // want "use assert\\.Empty"
assert.Emptyf(t, arr, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
		},
		{
			name: "full mode",
			mode: (*AssertionExpander).FullMode,
			expected: `assert.Len(t, arr, 0) // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg") // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg with arg %d", 42) // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg with args %d %s", 42, "42") // want "use assert\\.Empty"
assert.Lenf(t, arr, 0, "msg") // want "use assert\\.Emptyf"
assert.Lenf(t, arr, 0, "msg with arg %d", 42) // want "use assert\\.Emptyf"
assert.Lenf(t, arr, 0, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
			expectedGolden: `assert.Empty(t, arr) // want "use assert\\.Empty"
assert.Empty(t, arr, "msg") // want "use assert\\.Empty"
assert.Empty(t, arr, "msg with arg %d", 42) // want "use assert\\.Empty"
assert.Empty(t, arr, "msg with args %d %s", 42, "42") // want "use assert\\.Empty"
assert.Emptyf(t, arr, "msg") // want "use assert\\.Emptyf"
assert.Emptyf(t, arr, "msg with arg %d", 42) // want "use assert\\.Emptyf"
assert.Emptyf(t, arr, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
		},
		{
			name:           "fmt single mode",
			mode:           (*AssertionExpander).FmtSingleMode,
			expected:       `assert.Lenf(t, arr, 0, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
			expectedGolden: `assert.Emptyf(t, arr, "msg with args %d %s", 42, "42") // want "use assert\\.Emptyf"`,
		},
		{
			name: "not fmt set mode",
			mode: (*AssertionExpander).NotFmtSetMode,
			expected: `assert.Len(t, arr, 0) // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg") // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg with arg %d", 42) // want "use assert\\.Empty"
assert.Len(t, arr, 0, "msg with args %d %s", 42, "42") // want "use assert\\.Empty"`,
			expectedGolden: `assert.Empty(t, arr) // want "use assert\\.Empty"
assert.Empty(t, arr, "msg") // want "use assert\\.Empty"
assert.Empty(t, arr, "msg with arg %d", 42) // want "use assert\\.Empty"
assert.Empty(t, arr, "msg with args %d %s", 42, "42") // want "use assert\\.Empty"`,
		},
		{
			name:           "not fmt single mode",
			mode:           (*AssertionExpander).NotFmtSingleMode,
			expected:       `assert.Len(t, arr, 0) // want "use assert\\.Empty"`,
			expectedGolden: `assert.Empty(t, arr) // want "use assert\\.Empty"`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			newExpander := func() *AssertionExpander {
				expander := NewAssertionExpander()
				if tt.mode != nil {
					expander = tt.mode(expander)
				}
				return expander
			}

			t.Run("errored", func(t *testing.T) {
				got := newExpander().Expand(assrn, selector, "t", nil)
				assertEqualStrings(t, tt.expected, got)
			})

			t.Run("golden", func(t *testing.T) {
				got := newExpander().AsGolden().Expand(assrn, selector, "t", nil)
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
