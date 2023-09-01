# testifylint

[![CI](https://github.com/Antonboom/testifylint/actions/workflows/ci.yml/badge.svg)](https://github.com/Antonboom/testifylint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Antonboom/testifylint)](https://goreportcard.com/report/github.com/Antonboom/testifylint)
[![Coverage](https://coveralls.io/repos/github/Antonboom/testifylint/badge.svg?branch=master)](https://coveralls.io/github/Antonboom/testifylint?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Antonboom/testifylint/blob/master/CONTRIBUTING.md#Open-for-contribution)

Checks usage of [github.com/stretchr/testify](https://github.com/stretchr/testify).

## Problem statement

Tests are also program code and the requirements for them should not differ much from the requirements
for the code under tests 🙂

We should try to maintain the consistency of tests, increase their readability, reduce the chance of bugs
and speed up the search for a problem.

[testify](https://github.com/stretchr/testify) is the most popular Golang testing framework* in recent years.
But it has a terrible ambiguous API in places, and **the purpose of this linter is to protect you from annoying mistakes**.

Most checkers are stylistic, but checkers like [error-is](#error-is), [require-error](#require-error),
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

```
$ testifylint --enable-all ./...
$ testifylint --enable=empty,error-is ./...
$ testifylint --enable=expected-actual --expected-actual.pattern=^wanted$ ./...
$ testifylint --enable=suite-extra-assert-call --suite-extra-assert-call.mode=require ./...
```

## Checkers

| Name                                                | Enabled By Default | Autofix |
|-----------------------------------------------------|--------------------|---------|
| [bool-compare](#bool-compare)                       | ✅                  | ✅       |
| [compares](#compares)                               | ✅                  | ✅       |
| [empty](#empty)                                     | ✅                  | ✅       |
| [error-is](#error-is)                               | ✅                  | ✅       |
| [error-nil](#error-nil)                             | ✅                  | ✅       |
| [expected-actual](#expected-actual)                 | ✅                  | ✅       |
| [float-compare](#float-compare)                     | ✅                  | ❌       |
| [len](#len)                                         | ✅                  | ✅       |
| [require-error](#require-error)                     | ✅                  | ❌       |
| [suite-dont-use-pkg](#suite-dont-use-pkg)           | ✅                  | ✅       |
| [suite-extra-assert-call](#suite-extra-assert-call) | ❌                  | ✅       |
| [suite-thelper](#suite-thelper)                     | ❌                  | ✅       |

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
     // And other variations around len(arr)...

✅   assert.Empty(t, arr)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### error-is

```go
❌   assert.Error(t, err, errSentinel) // Typo, errSentinel hits `msgAndArgs`.
     assert.NoError(t, err, errSentinel)
     assert.True(t, errors.Is(err, errSentinel))
     assert.False(t, errors.Is(err, errSentinel))

✅   assert.ErrorIs(t, err, errSentinel)
     assert.NotErrorIs(t, err, errSentinel)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: In the first two cases, a common mistake that leads to hiding the incorrect wrapping of sentinel errors.

---

### error-nil

```go
❌   assert.Nil(t, err)
     assert.NotNil(t, err)
     assert.Equal(t, err, nil)
     assert.NotEqual(t, err, nil)

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
(`-expected-actual.pattern` flag).

It is planned [to change the order of assertion arguments](https://github.com/stretchr/testify/issues/1089#Argument_order) to more natural
(actual, expected) in `v2`.

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

### len

```go
❌   assert.Equal(t, 3, len(arr))
     assert.True(t, len(arr) == 5)

✅   assert.Len(t, arr, 3)
     assert.Len(t, arr, 5)
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
**Reason**: "Ignoring" errors is not the "Go way" and leads to further panics in the test, making it harder to debug.

`tesitfy/require` allows to stop test execution when a test fails.

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

You can enable such behaviaor through `--suite-extra-assert-call.mode=require`.

**Autofix**: true. <br>
**Enabled by default**: true, in `remove` mode. <br>
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
**Reason**: Consistency to non-suite test helpers. Explicit markup of helper methods.

`s.T().Helper()` call is not important technically because `testify` prints full `Error Trace`
[anyway](https://github.com/stretchr/testify/blob/882382d845cd9780bd93c1acc8e1fa2ffe266ca1/assert/assertions.go#L317).

The checker rather acts as an example of a [checkers.AdvancedChecker](https://github.com/Antonboom/testifylint/blob/676324836555445fded4e9afc004101ec6f597fe/internal/checkers/checker.go#L56).

---

## testify V2

The second version of `testify` [promises](https://github.com/stretchr/testify/issues/1089) more "pleasant" API and
makes some above checkers irrelevant.

It will not be difficult to add `v2` support in the linter in the future.

But at the moment it looks like we are [extremely far](https://github.com/stretchr/testify/pull/1109#issuecomment-1650619745)
from `v2`. Related milestone [here](https://github.com/stretchr/testify/milestone/4).
