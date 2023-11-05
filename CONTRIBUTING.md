# Contribution guideline

> ⚠️ The module version is **1.20**, but for development you need Go >= **1.21**

### 1) Install developer tools

```bash
# https://taskfile.dev/installation/
$ go install github.com/go-task/task/v3/cmd/task@latest
$ task tools:install
Install dev tools...
```

### 2) Create a checker skeleton in `internal/checkers/{checker_name}.go`

For example, we want to create a new checker
`Zero` in `internal/checkers/zero.go`:

```go
package checkers

// Zero detects situations like
//
//	assert.Equal(t, 0, count) 
//	assert.Equal(t, nil, userObj)
//
// and requires
//
//	assert.Zero(t, count) 
//	assert.Zero(t, userObj)
type Zero struct{}

// NewZero constructs Zero checker.
func NewZero() Zero      { return Zero{} }
func (Zero) Name() string { return "zero" }
```

The above code is enough to satisfy the `checkers.Checker` interface.

### 3) Add the new checker to the `registry` in order of priority

The earlier the checker is in [the registry](internal/checkers/checkers_registry.go), the more priority it is.

For example, the `zero` checker takes precedence over the `expected-actual` or `empty`, 
because its check is more "narrow" and when you fix the warning from `zero`, 
the rest of the checkers will become irrelevant.

```go
var registry = checkersRegistry{
    // ...
    {factory: asCheckerFactory(NewZero), enabledByDefault: false},
    // ...
    {factory: asCheckerFactory(NewEmpty), enabledByDefault: true},
    // ...
    {factory: asCheckerFactory(NewExpectedActual), enabledByDefault: true},
    // ...
}
```

By default, we disable the checker if we doubt its 100% usefulness.

### 4) Create new tests generator in `internal/testgen/gen_{checker_name}.go`

Create new `ZeroTestsGenerator` in `internal/testgen/gen_zero.go`.

See examples in adjacent files.

In the first iteration, these can be a very simple tests for debugging checker's proof of concept.
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

`Zero` is an example of [checkers.RegularChecker](./internal/checkers/checker.go) because it works with "general" 
assertion call. For more complex checkers, use the [checkers.AdvancedChecker](./internal/checkers/checker.go) interface.

If the checker turns out to be too “fat”, then you can omit some obviously rare combinations,
especially if they are covered by other checkers. Usually these are expressions in `assert.True/False`.

### 8) Improve tests from p.4 if necessary

Pay attention to `Assertion` and `NewAssertionExpander`, but keep your tests as small as possible.
Usually 100-200 lines are enough for checker testing. You need to find balance between coverage,
common sense, and processing of boundary conditions. See existing tests as example.

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
- [inefficient-assert](#inefficient-assert)
- [negative-positive](#negative-positive)
- [no-fmt-mess](#no-fmt-mess)
- [require-len](#require-len)
- [suite-run](#suite-run)
- [suite-test-name](#suite-test-name)
- [zero](#zero)

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
     // etc.

✅   assert.ErrorIs(t, err, ErrUserNotFound)
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: The `Error()` method on the `error` interface exists for humans, not code.

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
**Reason**: Possible resource leaks, because `require` finishes the current goroutine.

---

### equal-values

```go
❌   assert.Equal(t, int64(100), price.Amount)
✅   assert.EqualValues(t, 100, price.Amount)
```

**Autofix**: true. <br>
**Enabled by default**: maybe? <br>
**Reason**: Code simplification.

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
**Reason**: Is similar to the [usestdlibvars](https://golangci-lint.run/usage/linters/#usestdlibvars) linter.

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
**Reason**: Code simplification.

---

### inefficient-assert

Simple:
```go
body, err := io.ReadAll(rr.Body)
require.NoError(t, err)
require.NoError(t, err) ❌

expectedJSON, err := json.Marshal(expected)
require.NoError(t, err)
require.JSONEq(t, string(expectedJSON), string(body))
```

Complex:
```go
require.NoError(t, err)
assert.ErrorContains(t, err, "user") ❌
```

**Autofix**: false. <br>
**Enabled by default**: Probably false, depends on implementation performance. <br>
**Reason**: Code simplification, elimination of possible bugs.

---

### negative-positive

```go
❌   assert.Less(t, val, 0)
     assert.Greater(t, val, 0)

✅   assert.Negative(t, val)
     assert.Positive(t, val)
```

**Autofix**: true. <br>
**Enabled by default**: maybe? <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### no-fmt-mess

**Autofix**: true. <br>
**Enabled by default**: maybe?

Those who are new to `testify` may be discouraged by the duplicative API:

```go
func Equal(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
func Equalf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool
```

The f-functions [were added a long time ago](https://github.com/stretchr/testify/pull/356) to eliminate 
`govet` [complain](https://go.googlesource.com/tools/+/refs/heads/release-branch.go1.12/go/analysis/passes/printf/printf.go?pli=1#506).

This introduces some inconsistency into the test code, and the next strategies are seen for the checker:

1) Forbid f-functions at all (also could be done through the [forbidigo](https://golangci-lint.run/usage/linters/#forbidigo) linter).

This will make it easier to migrate to [v2](https://github.com/stretchr/testify/issues/1089), because

> Format functions should not be accepted as they are equivalent to their "non-f" counterparts.

But it doesn't look like a "go way" and the `govet` won't be happy.

2) IMHO, a more appropriate option is to disallow the use of `msgAndArgs` in non-f assertions. Look at 
[the comment](https://github.com/stretchr/testify/issues/1089#issuecomment-1695059265).

But there will be no non-stylistic benefits from the checker in this case (depends on the view of API in `v2`).

Also, in the first iteration `no-fmt-mess` checker could help to avoid code like this:
```go
assert.Error(t, err, fmt.Sprintf("Profile %s should not be valid", test.profile))
```

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

### suite-run

```go
func (s *Suite) TestSomething() {
    ❌
    s.T().Run("subtest1", func(t *testing.T) {
        // ...
        assert.Equal(t, "gopher", result)
    })
	
    ✅
    s.Run("subtest1", func() {
        // ...
        s.Equal("gopher", result)
    })
}
```

**Autofix**: true. <br>
**Enabled by default**: probably yes. <br>
**Reason**: Code simplification and consistency.

But need to investigate the technical difference and the reasons for the appearance of `s.Run`.
Also, maybe this case is already covered by [suite-dont-use-pkg](README.md#suite-dont-use-pkg)?

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
**Reason**: Just unification of approach.

Also, maybe to check the configurable format of subtest name? Mess example:

```go
func (s *HandlersSuite) Test_Usecase_Success()
func (s *HandlersSuite) TestUsecaseSuccess()
func (s *HandlersSuite) Test_UsecaseSuccess()
```

---

### zero

```go
❌   assert.Equal(t, 0, count)
     assert.Equal(t, nil, userObj)
     assert.Equal(t, "", name)
     // etc.

✅   assert.Zero(t, count)
     assert.Zero(t, userObj)
     assert.Zero(t, name)
```

**Autofix**: true. <br>
**Enabled by default**: false. <br>
**Reason**: Just for your reflection and suggestion.

I'm not sure if anyone uses `assert.Zero` – it looks strange and conflicts with `assert.Empty`:

```go
❌   assert.Equal(t, "", result)
     assert.Nil(t, errCh)

✅   assert.Empty(t, result)
     assert.Empty(t, errCh)
```

Maybe it's better to make configurable support for other types in the `empty` checker and
vice versa to prohibit the `Zero`?

---

Any other figments of your imagination are welcome 🙏<br>
List of `testify` functions [here](https://pkg.go.dev/github.com/stretchr/testify@master/assert#pkg-functions).
