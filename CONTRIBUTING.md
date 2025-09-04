# Contribution guideline

### 1) Install `task`

```bash
# https://taskfile.dev/installation/
$ go install github.com/go-task/task/v3/cmd/task@latest
```

### 2) Create a checker skeleton in `internal/checkers/{checker_name}.go`

For example, we want to create a new checker
`TimeCompare` in `internal/checkers/time_compare.go`:

```go
package checkers

// TimeCompare detects situations like
//
//	assert.Equal(t, expTs, actualTs)
//
// and requires
//
//	assert.True(t, actualTs.Equal(expTs))
type TimeCompare struct{}

// NewTimeCompare constructs TimeCompare checker.
func NewTimeCompare() TimeCompare { return TimeCompare{} }
func (TimeCompare) Name() string  { return "time-compare" }
```

The above code is enough to satisfy the `checkers.Checker` interface.

### 3) Add the new checker to the `registry` in order of priority

The earlier the checker is in [the registry](internal/checkers/checkers_registry.go), the more priority it is.

For example, the `TimeCompare` checker takes precedence over the `empty` and `expected-actual`,
because its check is more "narrow" and when you fix the warning from `TimeCompare`,
the rest of the checkers will become irrelevant.

```go
var registry = checkersRegistry{
    // ...
    {factory: asCheckerFactory(NewTimeCompare), enabledByDefault: true},
    // ...
    {factory: asCheckerFactory(NewEmpty), enabledByDefault: true},
    // ...
    {factory: asCheckerFactory(NewExpectedActual), enabledByDefault: true},
    // ...
}
```

By default, we disable the checker if we doubt its 100% usefulness.

### 4) Create new tests generator in `internal/testgen/gen_{checker_name}.go`

Create new `TimeCompareTestsGenerator` in `internal/testgen/gen_time_compare.go`.

See examples in adjacent files, e.g. [internal/testgen/gen_regexp.go](internal/testgen/gen_regexp.go).

In the first iteration, these can be a simple tests for debugging checker's proof of concept.
And after the implementation of the checker, you can add various cycles, variables, etc. to the template.

`GoldenTemplate` is usually an `ErroredTemplate` with some strings replaced.

### 5) Add generator into `checkerTestsGenerators`

Look at [internal/testgen/main.go](./internal/testgen/main.go).

### 6) Generate and run tests

Tests should fail.

```bash
$ task test
Generate analyzer tests...
Test...
...
--- FAIL: TestTestifyLint_CheckersDefault
FAIL
```

### 7) Implement the checker

`TimeCompare` is an example of [checkers.RegularChecker](./internal/checkers/checker.go) because it works with "general"
assertion call. For more complex checkers, use the [checkers.AdvancedChecker](./internal/checkers/checker.go) interface.

If the checker turns out to be too “fat”, then you can omit some obviously rare combinations,
especially if they are covered by other checkers. Usually these are expressions in `assert.True/False`.

Remember that [assert.TestingT](https://pkg.go.dev/github.com/stretchr/testify/assert#TestingT) and
[require.TestingT](https://pkg.go.dev/github.com/stretchr/testify/require#TestingT) are different interfaces,
which may be important in some contexts.

Also, pay attention to `internal/checkers/helpers_*.go` files. Try to reuse existing code as much as possible.

### 8) Improve tests from p.4 if necessary

Pay attention to `Assertion` and `NewAssertionExpander`, but keep your tests as small as possible.
Usually 100-200 lines are enough for checker testing. You need to find balance between coverage,
common sense, and processing of boundary conditions. See existing tests as example.

For testing checker replacements use [testdata/src/debug](./analyzer/testdata/src/debug).

Example of test cases search – https://github.com/search?q=language%3Ago+%2Fassert%5C.IsType%5C%28t%5C%2C.*err%2F&type=code

### 9) Run the `task` command from the project's root directory

```bash
$ task
Tidy...
Fmt...
Lint...
Generate analyzer tests...
Test...
Install...
```

Fix linter issues and broken tests (probably related to the checkers registry).

To run checker default tests you can use `task test:checker -- {checker-name}`, e.g.

```bash
$ task test:checker -- time-compare
```

### 10) Update `README.md`, commit the changes and submit a pull request 🔥

Describe a new checker in [checkers section](./README.md#checkers).

# Open for contribution

- [elements-match](#elements-match)
- [error-compare](#error-compare)
- [equal-values](#equal-values)
- [graceful-teardown](#graceful-teardown)
- [float-compare](#float-compare)
- [http-const](#http-const)
- [http-sugar](#http-sugar)
- [require-len](#require-len)
- [suite-test-name](#suite-test-name)

---

### elements-match

```go
❌   require.Equal(t, len(expected), len(result)
     sort.Slice(expected, /* ... */)
     sort.Slice(result, /* ... */)
     for i := range result {
         assert.Equal(t, expected[i], result[i])
     }
     // Or for Go >= 1.21
     slices.Sort(expected)
     slices.Sort(result)
     assert.True(t, slices.Equal(expected, result))

✅   assert.ElementsMatch(t, expected, result)
```

**Autofix**: maybe (depends on implementation difficulty). <br>
**Enabled by default**: maybe (depends on checker's stability). <br>
**Reason**: Code simplification.

---

### error-compare

```go
❌   assert.ErrorContains(t, err, "not found")
     assert.EqualError(t, err, "user not found")
     assert.Equal(t, err.Error(), "user not found")
     assert.Equal(t, err, errSentinel) // Through `reflect.DeepEqual` causes error strings to be compared.
     assert.NotEqual(t, err, errSentinel)
     require.Error(t, fmt.Errorf("you need to specify either logGroupName or logGroupArn"), err) // grafana case
     // etc.

✅   assert.ErrorIs(t, err, ErrUserNotFound)
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: The `Error()` method on the `error` interface exists for humans, not code. <br>
**Related issues**: [#47](https://github.com/Antonboom/testifylint/issues/47)

---

### graceful-teardown

Warns about usage of `require` in `t.Cleanup` functions and suite teardown methods:

```go
func (s *ServiceIntegrationSuite) TearDownTest() {
    if p := s.verdictsProducer; p != nil {
        s.Require().NoError(p.Close()) ❌
    }
    if c := s.dlqVerdictsConsumer; c != nil {
        s.NoError(c.Close())
    }
    s.DBSuite.TearDownTest()
    s.ks.TearDownTest()
}
```

**Autofix**: false. <br>
**Enabled by default**: false. <br>
**Reason**: Possible resource leaks, because `require` finishes the current goroutine. <br>
**Related issues**: [#142](https://github.com/Antonboom/testifylint/issues/142)

---

### float-compare

1) Support other tricky cases

```go
❌   assert.NotEqual(t, 42.42, a)
     assert.Greater(t, a, 42.42)
     assert.GreaterOrEqual(t, a, 42.42)
     assert.Less(t, a, 42.42)
     assert.LessOrEqual(t, a, 42.42)

     assert.True(t, a != 42.42) // assert.False(t, a == 42.42)
     assert.True(t, a > 42.42)  // assert.False(t, a <= 42.42)
     assert.True(t, a >= 42.42) // assert.False(t, a < 42.42)
     assert.True(t, a < 42.42)  // assert.False(t, a >= 42.42)
     assert.True(t, a <= 42.42) // assert.False(t, a > 42.42)
```

But I have no idea for equivalent. Probably we need a new API from `testify`.

2) Support compares of "float containers" (structs, slices, arrays, maps, something else?), e.g.

```go
type Tx struct {
    ID string
    Score float64
}

❌   assert.Equal(t, Tx{ID: "xxx", Score: 0.9643}, tx)

✅   assert.Equal(t, "xxx", tx.ID)
     assert.InEpsilon(t, 0.9643, tx.Score, 0.0001)
```

And similar idea for `assert.InEpsilonSlice` / `assert.InDeltaSlice`.

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Work with floating properly.

---

### http-const

```go
❌   assert.HTTPStatusCode(t, handler, "GET", "/index", nil, 200)
     assert.HTTPBodyContains(t, handler, "GET", "/index", nil, "counter")
     // etc.

✅   assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, http.StatusOK)
     assert.HTTPBodyContains(t, handler, http.MethodGet, "/index", nil, "counter")
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: Is similar to the [usestdlibvars](https://golangci-lint.run/usage/linters/#usestdlibvars) linter. <br>
**Related issues**: [#141](https://github.com/Antonboom/testifylint/issues/141)

---

### http-sugar

```go
❌   assert.HTTPStatusCode(t,
         handler, http.MethodGet, "/index", nil, http.StatusOK)
     assert.HTTPStatusCode(t,
         handler, http.MethodGet, "/admin", nil, http.StatusNotFound)
     assert.HTTPStatusCode(t,
         handler, http.MethodGet, "/oauth", nil, http.StatusFound)
     // etc.

✅   assert.HTTPSuccess(t, handler, http.MethodGet, "/index", nil)
     assert.HTTPError(t, handler, http.MethodGet, "/admin", nil)
     assert.HTTPRedirect(t, handler, http.MethodGet, "/oauth", nil)
```

**Autofix**: true. <br>
**Enabled by default**: maybe? <br>
**Reason**: Code simplification. <br>
**Related issues**: [#140](https://github.com/Antonboom/testifylint/issues/140)

---

### require-len

The main idea: if code contains array/slice indexing,
then before that there must be a length constraint through `require`.

```go
❌   assert.Len(t, arr, 3) // Or assert.NotEmpty(t, arr) and other variations.
     assert.Equal(t, "gopher", arr[1])

✅   require.Len(t, arr, 3) // Or require.NotEmpty(t, arr) and other variations.
     assert.Equal(t, "gopher", arr[1])
```

**Autofix**: false? <br>
**Enabled by default**: maybe? <br>
**Reason**: Similar to [require-error](README.md#require-error). Save you from annoying panics.

Or maybe do something similar for maps? And come up with better name for the checker.

---

### suite-test-name

```go
import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type BalanceSubscriptionSuite struct {
    suite.Suite
}

❌ func TestBalanceSubs_Run(t *testing.T) {
    suite.Run(t, new(BalanceSubscriptionSuite))
}


✅ func TestBalanceSubscriptionSuite(t *testing.T) {
    suite.Run(t, new(BalanceSubscriptionSuite))
}
```

**Autofix**: true. <br>
**Enabled by default**: false. <br>
**Reason**: Just unification of approach. <br>
**Related issues**: [#48](https://github.com/Antonboom/testifylint/issues/48)

Also, maybe to check the configurable format of subtest name? Mess example:

```go
func (s *HandlersSuite) Test_Usecase_Success()
func (s *HandlersSuite) TestUsecaseSuccess()
func (s *HandlersSuite) Test_UsecaseSuccess()
```

---

Any other figments of your imagination are welcome 🙏<br>
List of `testify` functions [here](https://pkg.go.dev/github.com/stretchr/testify@master/assert#pkg-functions).

# FAQ

### Why do we use `internal/testify` instead of `github.com/stretchr/testify`?

1) [internal/testify](./internal/testify) is not a local copy of `stretchr/testify`. The package contains domain (for 
linter context) entities, which absent in `testify` itself.

2) We cannot depend on `stretch/testify`, because it causes dependency bomb for linter's module. We should keep as min
   dependencies as possible.

### Why do we import `github.com/stretchr/testify` in tests?

Such imports in `testdata` do not affect linter module's dependencies. Moreover, of course, tests should use the real
`testify`.
