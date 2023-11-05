# testifylint

![Latest release](https://img.shields.io/github/v/release/Antonboom/testifylint)
[![CI](https://github.com/Antonboom/testifylint/actions/workflows/ci.yml/badge.svg)](https://github.com/Antonboom/testifylint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Antonboom/testifylint)](https://goreportcard.com/report/github.com/Antonboom/testifylint?dummy=unused)
[![Coverage Status](https://coveralls.io/repos/github/Antonboom/testifylint/badge.svg?branch=master)](https://coveralls.io/github/Antonboom/testifylint?branch=master&dummy=unused)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Antonboom/testifylint/blob/master/CONTRIBUTING.md#Open-for-contribution)

Checks usage of [github.com/stretchr/testify](https://github.com/stretchr/testify).

## Table of Contents

* [Problem statement](#problem-statement)
* [Installation & usage](#installation--usage)
* [Configuring](#configuring)
* [Checkers](#checkers)
* [Chain of warnings](#chain-of-warnings)
* [testify v2](#testify-v2)

## Problem statement

Tests are also program code and the requirements for them should not differ much from the requirements
for the code under tests 🙂

We should try to maintain the consistency of tests, increase their readability, reduce the chance of bugs
and speed up the search for a problem.

[testify](https://github.com/stretchr/testify) is the most popular Golang testing framework* in recent years.
But it has a terrible ambiguous API in places, and **the purpose of this linter is to protect you from annoying mistakes**.

Most checkers are stylistic, but checkers like [error-is-as](#error-is-as), [require-error](#require-error),
[expected-actual](#expected-actual), [float-compare](#float-compare) are really helpful.

_* JetBrains "The State of Go Ecosystem" reports [2021](https://www.jetbrains.com/lp/devecosystem-2021/go/#Go_which-testing-frameworks-do-you-use-regularly-if-any)
and [2022](https://www.jetbrains.com/lp/devecosystem-2022/go/#which-testing-frameworks-do-you-use-regularly-if-any-)._

## Installation & usage

```
$ go install github.com/Antonboom/testifylint@latest
$ testifylint -h
$ testifylint ./...
```

## Configuring

### CLI
```bash
# Enable all checkers.
$ testifylint --enable-all ./...

# Enable specific checkers only.
$ testifylint --disable-all --enable=empty,error-is-as ./...

# Disable specific checkers only.
$ testifylint --enable-all --disable=empty,error-is-as ./...

# Checkers configuration.
$ testifylint --expected-actual.pattern=^wanted$ ./...
$ testifylint --require-error.fn-pattern="^(Errorf?|NoErrorf?)$" ./...
$ testifylint --suite-extra-assert-call.mode=require ./...
```

### golangci-lint
https://golangci-lint.run/usage/linters/#testifylint

## Checkers

| Name                                                | Enabled By Default | Autofix |
|-----------------------------------------------------|--------------------|---------|
| [bool-compare](#bool-compare)                       | ✅                  | ✅       |
| [compares](#compares)                               | ✅                  | ✅       |
| [empty](#empty)                                     | ✅                  | ✅       |
| [error-is-as](#error-is-as)                         | ✅                  | ✅       |
| [error-nil](#error-nil)                             | ✅                  | ✅       |
| [expected-actual](#expected-actual)                 | ✅                  | ✅       |
| [float-compare](#float-compare)                     | ✅                  | ❌       |
| [go-require](#go-require)                           | ✅                  | ❌       |
| [len](#len)                                         | ✅                  | ✅       |
| [nil-compare](#nil-compare)                         | ✅                  | ✅       |
| [require-error](#require-error)                     | ✅                  | ❌       |
| [suite-dont-use-pkg](#suite-dont-use-pkg)           | ✅                  | ✅       |
| [suite-extra-assert-call](#suite-extra-assert-call) | ✅                  | ✅       |
| [suite-thelper](#suite-thelper)                     | ❌                  | ✅       |

> ⚠️ Also look at open for contribution [checkers](CONTRIBUTING.md#open-for-contribution)

---

### bool-compare

```go
❌   assert.Equal(t, true, result)
     assert.NotEqual(t, true, result)
     assert.False(t, !result)
     assert.True(t, result == true)
     // And other variations...

✅   assert.True(t, result)
     assert.False(t, result)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: Code simplification.

---

### compares

```go
❌   assert.True(t, a == b)
     assert.True(t, a != b)
     assert.True(t, a > b)
     assert.True(t, a >= b)
     assert.True(t, a < b)
     assert.True(t, a <= b)
     // And other variations (with assert.False too)...

✅   assert.Equal(t, a, b)
     assert.NotEqual(t, a, b)
     assert.Greater(t, a, b)
     assert.GreaterOrEqual(t, a, b)
     assert.Less(t, a, b)
     assert.LessOrEqual(t, a, b)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### empty

```go
❌   assert.Len(t, arr, 0)
     assert.Equal(t, 0, len(arr))
     assert.NotEqual(t, 0, len(arr))
     assert.GreaterOrEqual(t, len(arr), 1)
     // And other variations around len(arr)...

✅   assert.Empty(t, arr)
     assert.NotEmpty(t, err)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### error-is-as

```go
❌   assert.Error(t, err, errSentinel) // Typo, errSentinel hits `msgAndArgs`.
     assert.NoError(t, err, errSentinel)
     assert.True(t, errors.Is(err, errSentinel))
     assert.False(t, errors.Is(err, errSentinel))
     assert.True(t, errors.As(err, &target))

✅   assert.ErrorIs(t, err, errSentinel)
     assert.NotErrorIs(t, err, errSentinel)
     assert.ErrorAs(t, err, &target)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: In the first two cases, a common mistake that leads to hiding the incorrect wrapping of sentinel errors.
In the rest cases – more appropriate `testify` API with clearer failure message.

Also `error-is-as` repeats `go vet`'s [errorsas check](https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/errorsas/errorsas.go) 
logic, but without autofix.

---

### error-nil

```go
❌   assert.Nil(t, err)
     assert.NotNil(t, err)
     assert.Equal(t, err, nil)
     assert.NotEqual(t, err, nil)
     assert.ErrorIs(t, err, nil)
     assert.NotErrorIs(t, err, nil)

✅   assert.NoError(t, err)
     assert.Error(t, err)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

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
**Enabled by default**: true. <br>
**Reason**: A common mistake that makes it harder to understand the reason of failed test.

The checker considers the expected value to be a basic literal, constant, or variable whose name matches the pattern
(`--expected-actual.pattern` flag).

It is planned [to change the order of assertion arguments](https://github.com/stretchr/testify/issues/1089#Argument_order) to more natural
(actual, expected) in `v2` of `testify`.

---

### float-compare

```go
❌   assert.Equal(t, 42.42, a)
     assert.True(t, a == 42.42)
     assert.False(t, a != 42.42)
	
✅   assert.InEpsilon(t, 42.42, a, 0.0001)
     assert.InDelta(t, 42.42, a, 0.01)
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Do not forget about [floating point rounding issues](https://floating-point-gui.de/errors/comparison/).

This checker is similar to the [floatcompare](https://github.com/golangci/golangci-lint/pull/2608) linter.

---

### go-require

```go
go func() {
    conn, err = lis.Accept()
    require.NoError(t, err) ❌
    
    if assert.Error(err) {
        assert.FailNow(t, msg) ❌
    }
}()
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Incorrect use of functions.

This checker is a radically improved analogue of `go vet`'s
[testinggoroutine](https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/testinggoroutine/doc.go) check.

The point of the check is that, according to the [documentation](https://pkg.go.dev/testing#T), 
functions leading to `t.FailNow` (essentially to `runtime.GoExit`) must only be used in the goroutine that runs the test.
Otherwise, they will not work as declared, namely, finish the test function.

You can disable the `go-require` checker and continue to use `require` as the current goroutine finisher, but this could lead
1. to possible resource leaks in tests;
2. to increasing of confusion, because functions will be not used as intended.

Typically, any assertions inside goroutines are a marker of poor test architecture.
Try to execute them in the main goroutine and distribute the data necessary for this into it
([example](https://github.com/ipfs/kubo/issues/2043#issuecomment-164136026)).

Also a bad solution would be to simply replace all `require` in goroutines with `assert`
(like [here](https://github.com/gravitational/teleport/pull/22567/files#diff-9f5fd20913c5fe80c85263153fa9a0b28dbd1407e53da4ab5d09e13d2774c5dbR7377))
– this will only mask the problem.

The checker is enabled by default, because `testinggoroutine` is enabled by default in `go vet`.

P.S. Related `testify`'s [thread](https://github.com/stretchr/testify/issues/772).

---

### len

```go
❌   assert.Equal(t, 3, len(arr))
     assert.True(t, len(arr) == 3)

✅   assert.Len(t, arr, 3)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### nil-compare

```go
❌   assert.Equal(t, value, nil)
     assert.NotEqual(t, value, nil)

✅   assert.Nil(t, value)
     assert.NotNil(t, value)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### require-error

```go
❌   assert.NoError(t, err)
     s.ErrorIs(err, io.EOF)
     s.Assert().Error(err)
     // And other error assertions...

✅   require.NoError(t, err)
     s.Require().ErrorIs(err, io.EOF)
     s.Require().Error(err)
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Such "ignoring" of errors leads to further panics, making the test harder to debug.

[testify/require](https://pkg.go.dev/github.com/stretchr/testify@master/require#hdr-Assertions) allows 
to stop test execution when a test fails.

To minimize the number of false positives, `require-error` ignores:
- assertion in the `if` condition;
- the entire `if-else` block, if there is an assertion in the `if` condition;
- the last assertion in the block, if there are no methods/functions calls after it;
- assertions in an explicit goroutine;
- assertions in an explicit testing cleanup function or suite teardown methods;
- sequence of `NoError` assertions.

Also you can configure functions to analyze via `--require-error.fn-pattern` flag.

---

### suite-dont-use-pkg

```go
import "github.com/stretchr/testify/assert"

func (s *MySuite) TestSomething() {
     ❌ assert.Equal(s.T(), 42, value)
     ✅ s.Equal(42, value)
}
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More simple and uniform code.

---

### suite-extra-assert-call

By default, the checker wants you to remove unnecessary `Assert()` calls:

```go
func (s *MySuite) TestSomething() {
     ❌ s.Assert().Equal(42, value)
     ✅ s.Equal(42, value)
}
```

But sometimes, on the contrary, people want consistency with `s.Assert()` and `s.Require()`:

```go
func (s *MySuite) TestSomething() {
     // ...
 
     ❌
     s.Require().NoError(err)
     s.Equal(42, value)

     ✅
     s.Require().NoError(err)
     s.Assert().Equal(42, value)
}
```

You can enable such behavior through `--suite-extra-assert-call.mode=require`.

**Autofix**: true. <br>
**Enabled by default**: true, in the `remove` mode. <br>
**Reason**: More simple or uniform code.

---

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
**Enabled by default**: false. <br>
**Reason**: Consistency with non-suite test helpers. Explicit markup of helper methods.

`s.T().Helper()` call is not important actually because `testify` prints full `Error Trace`
[anyway](https://github.com/stretchr/testify/blob/882382d845cd9780bd93c1acc8e1fa2ffe266ca1/assert/assertions.go#L317).

The checker rather acts as an example of a [checkers.AdvancedChecker](https://github.com/Antonboom/testifylint/blob/676324836555445fded4e9afc004101ec6f597fe/internal/checkers/checker.go#L56).

---

## Chain of warnings

Linter does not automatically handle the "evolution" of changes.
And in some cases may be dissatisfied with your code several times, for example:

```go
assert.True(err == nil)   // compares: use assert.Equal
assert.Equal(t, err, nil) // error-nil: use assert.NoError
assert.NoError(t, err)    // require-error: for error assertions use require
require.NoError(t, err)
```

Please [contribute](./CONTRIBUTING.md) if you have ideas on how to make this better.

## testify v2

The second version of `testify` [promises](https://github.com/stretchr/testify/issues/1089) more "pleasant" API and
makes some above checkers irrelevant.

In this case, the possibility of supporting `v2` in the linter is not excluded.

But at the moment it looks like we are [extremely far](https://github.com/stretchr/testify/pull/1109#issuecomment-1650619745)
from `v2`. Related milestone [here](https://github.com/stretchr/testify/milestone/4).
