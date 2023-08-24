# testifylint
Checks usage of [github.com/stretchr/testify](https://github.com/stretchr/testify).

## Installation & usage

```
$ go install github.com/Antonboom/testifylint@latest
$ testifylint ./...
```

## Configuring
TODO: via golangci-lint.
- dump
- look at example

## Checkers

- [bool-compare](#bool-compare)
- [compares](#compares)
- [empty](#empty)
- [error](#error)
- [error-is](#error-is)
- [expected-actual](#expected-actual)
- [float-compare](#float-compare)
- [len](#len)
- [require-error](#require-error)
- [suite-dont-use-pkg](#suite-dont-use-pkg)
- [suite-no-extra-assert-call](#suite-no-extra-assert-call)
- [suite-thelper](#suite-thelper)

### bool-compare
```go
❌   assert.Equal(t, true, result)
     assert.Equal(t, false, result)
     assert.NotEqual(t, true, result)
     assert.False(t, !result)
    // And other variations...

✅   assert.True(t, result)
     assert.False(t, result)
     assert.False(t, result)
     assert.True(t, result)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### compares
```go
❌   assert.True(t, a == b)
     assert.True(t, a != b)
     assert.True(t, a > b)
     assert.True(t, a >= b)
     assert.True(t, a < b)
     assert.True(t, a <= b)
     // And other variations (assert.False including)...

✅   assert.Equal(t, a, b)
     assert.NotEqual(t, a, b)
     assert.Greater(t, a, b)
     assert.GreaterOrEqual(t, a, b)
     assert.Less(t, a, b)
     assert.LessOrEqual(t, a, b)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### empty
```go
❌   assert.Len(t, arr, 0)
     assert.Equal(t, len(arr), 0)
     // And other variations around len(arr)...

✅   assert.Empty(t, arr)
     assert.Empty(t, arr)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### error
```go
❌   assert.Nil(t, err)
     assert.NotNil(t, err)

✅   assert.NoError(t, err)
     assert.Error(t, err)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### error-is
```go
❌   assert.Error(t, err, errSentinel)
     assert.NoError(t, err, errSentinel)

✅   assert.ErrorIs(t, err, errSentinel)
     assert.NotErrorIs(t, err, errSentinel)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### expected-actual
```go
❌   assert.Equal(t, result, 42)
     assert.NotEqual(t, result, "expected")
     assert.JSONEq(t, result, `{"version": 3}`)
     assert.YAMLEq(t, result, "version: '3'")

✅   assert.Equal(t, 42, result)
     assert.NotEqual(t, "expected", result)
     assert.JSONEq(t, `{"version": 3}`, result)
     assert.YAMLEq(t, "version: '3'", result)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### float-compare
```go
❌   assert.Equal(t, 42.42, a)
     assert.True(t, a == 42.42)
     assert.False(t, a != 42.42)
	
✅   assert.InDelta(t, 42.42, a, 0.0001)
     assert.InEpsilon(t, 42.42, a, 0.01)
     assert.InDelta(t, 42.42, a, 0.001)
```
**Autofix**: false. <br>
**Enabled by default**: true.

### len
```go
❌   assert.Equal(t, len(arr), 3)
     assert.True(t, len(arr) == 5)

✅   assert.Len(t, arr, 3)
     assert.Len(t, arr, 5)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### require-error
```go
❌   assert.Error(t, err)
     assert.ErrorIs(t, err, io.EOF)
     assert.ErrorAs(t, err, new(os.PathError))
     assert.NoError(t, err)
     assert.NotErrorIs(t, err, io.EOF)
     s.Error(err)
     s.Assert().ErrorIs(err, io.EOF)
     // etc.

✅   require.Error(t, err)
     require.ErrorIs(t, err, io.EOF)
     require.ErrorAs(t, err, new(os.PathError))
     require.NoError(t, err)
     require.NotErrorIs(t, err, io.EOF)
     s.Require().Error(err)
     s.Require().ErrorIs(err, io.EOF)
```
**Autofix**: false. <br>
**Enabled by default**: true.

### suite-dont-use-pkg
```go
import "github.com/stretchr/testify/assert"

func (s *MySuite) TestSomething() {
    ❌ assert.Equal(s.T(), 42, value)
    ✅ s.Equal(42, value)
}
```
**Autofix**: true. <br>
**Enabled by default**: true.

### suite-no-extra-assert-call
```go
func (s *MySuite) TestSomething() {
    ❌ s.Assert().Equal(42, value)
    ✅ s.Equal(42, value)
}
```
**Autofix**: true. <br>
**Enabled by default**: false.

### suite-thelper
```go
❌
func (s *RoomSuite) assertRoomRound(roundID RoundID) {
    s.Equal(roundID, s.getRoom().CurrentRound.ID)
}

✅
func (s *RoomSuite) assertRoomRound(roundID RoundID) {
    s.T().Helper()
    s.Equal(roundID, s.getRoom().CurrentRound.ID)
}
```
**Autofix**: true. <br>
**Enabled by default**: false.

## Contribution
1) Create new test generator in [internal/cmd/testgen](internal/cmd/testgen/main.go).
If the checker strongly conflicts with the existing ones, then place the tests in a separate 
directory, otherwise in [testdata/src/checkers/most-of](pkg/analyzer/testdata/src/checkers/most-of).

2) Add new test in [pkg/analyzer/analyzer_test.go](pkg/analyzer/analyzer_test.go).

3) Implement new checker in [internal/checker](internal/checker).
Add the checker in `allCheckers` slice in [internal/checker/checkers.go](internal/checker/checkers.go).
Add the `DisabledByDefault()` method to the checker to disable it by default.
Set up the checker (if needed) in the [newCheckers](pkg/analyzer/checkers.go) function.

4) Add new checker in `Test_newCheckers`, `TestAllCheckers`, `TestEnabledByDefaultCheckers` and `TestDisabledByDefaultCheckers`.

5) Run the `task` command from the project's root directory.
```bash
$ cd testifylint
$ task
Install local tools...
Tidy...
Fmt...
Tests...
Install...
Dump config...
```

6) Update `README.md`, commit changes and submit a pull request 🙂.

<details>
  <summary>Open for contribution</summary>

### Global

### empty (existent checker)
Add config like
```yaml
empty:
  for-zero-valued-strings: true
  for-zero-valued-channels: false
  // ...
```

And support
```go
❌   assert.Equal(t, "", result)
     assert.Nil(t, errCh)

✅   assert.Empty(t, result)
     assert.Empty(t, errCh)
```

### float-compares (existent checker)
1) Support "structs with floats" compares
```go
type Tx struct {
    ID string
    Score float64
}

❌   assert.Equal(t, Tx{ID: "xxx", Score: 0.9643}, tx)

✅   assert.Equal(t, "xxx", tx.ID)
     assert.InDelta(t, 0.9643, tx.Score, 0.0001)
```

// TODO: slices
// require.Equal(t, [4]float64{0, .50, .50, 0}, rs.Column)

2) Support other cases
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
But I have no idea for equivalent. Probably we need a new API from testify.

## Proposed checkers
- [elements-match](#elements-match)
- [equal-values](#equal-values)
- [http-const](#http-const)
- [http-sugar](#http-sugar)
- [negative-positive](#negative-positive)
- [no-error-contains](#no-error-contains)
- [no-equal-error](#no-equal-error)
- [suite-test-name](#suite-test-name)
- [zero](#zero)

### elements-match
```go
❌   require.Equal(t, len(expected), len(result)
     sort.Slice(expected, /* ... */)
     sort.Slice(result, /* ... */)
     for i := range result {
         assert.Equal(t, expected[i], result[i])
     }
 
✅   assert.ElementsMatch(t, expected, result)
```
**Autofix**: maybe (depends on implementation difficulty). <br>
**Enabled by default**: maybe (depends on checker's stability).

### equal-values
```go
❌   assert.Equal(t, int64(100), price.Amount)
✅   assert.EqualValues(t, 100, price.Amount)
```
**Autofix**: true. <br>
**Enabled by default**: false.

### http-const
```go
❌   assert.HTTPStatusCode(t, handler, "GET", "/index", nil, 200)
     assert.HTTPBodyContains(t, handler, "GET", "/index", nil, "counter")
     // etc.

✅   assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, http.StatusOK)
     assert.HTTPBodyContains(t, handler, http.MethodGet, "/index", nil, "counter")
```
**Autofix**: true. <br>
**Enabled by default**: false.

### http-sugar
```go
❌   assert.HTTPStatusCode(t,
         handler, http.MethodGet, "/index", nil, http.StatusOK)
     assert.HTTPStatusCode(t, 
         handler, http.MethodGet, "/admin", nil, http.StatusNotFound)
     assert.HTTPStatusCode(t,
         handler, http.MethodGet, "/oauth", nil, http.StatusFound)

✅   assert.HTTPSuccess(t, handler, http.MethodGet, "/index", nil)
     assert.HTTPError(t, handler, http.MethodGet, "/admin", nil)
     assert.HTTPRedirect(t, handler, http.MethodGet, "/oauth", nil)
```
**Autofix**: true. <br>
**Enabled by default**: false.

### negative-positive
```go
❌   assert.Less(t, val, 0)
     assert.Greater(t, val, 0)

✅   assert.Negative(t, val)
     assert.Positive(t, val)
```
**Autofix**: true. <br>
**Enabled by default**: true.

### no-error-contains
```go
❌   assert.ErrorContains(t, err, "not found")
✅   assert.ErrorIs(t, err, ErrNotFound)
```
**Autofix**: false. <br>
**Enabled by default**: true.

### no-equal-error
```go
❌   assert.EqualError(t, err, "user not found")
✅   assert.ErrorIs(t, err, ErrUserNotFound)
```
**Autofix**: false. <br>
**Enabled by default**: true.

### suite-test-name
```go
import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type BalanceSubscriptionSuite struct {
    suite.Suite
}

❌
func TestBalanceSubs(t *testing.T) {
    suite.Run(t, new(BalanceSubscriptionSuite))
}

✅
func TestBalanceSubscriptionSuite(t *testing.T) {
    suite.Run(t, new(BalanceSubscriptionSuite))
}
```
**Autofix**: true. <br>
**Enabled by default**: false.

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
May conflict with the `empty` checker. <br>
**Autofix**: true. <br>
**Enabled by default**: false.

</details>
