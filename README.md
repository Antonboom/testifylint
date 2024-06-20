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
But it has a terrible ambiguous API in places, and **the purpose of this linter is to protect you from annoying
mistakes**.

Most checkers are stylistic, but checkers like [error-is-as](#error-is-as), [require-error](#require-error),
[expected-actual](#expected-actual), [float-compare](#float-compare) are really helpful.

_* JetBrains "The State of Go Ecosystem" reports
[2021](https://www.jetbrains.com/lp/devecosystem-2021/go/#Go_which-testing-frameworks-do-you-use-regularly-if-any)
and [2022](https://www.jetbrains.com/lp/devecosystem-2022/go/#which-testing-frameworks-do-you-use-regularly-if-any-)._

## Installation & usage

```
$ go install github.com/Antonboom/testifylint@latest
$ testifylint -h
$ testifylint ./...
```

### Fixing

```
$ testifylint --fix ./...
```

Fixing with `golangci-lint` is currently **unavailable** due to
[golangci/golangci-lint#1779](https://github.com/golangci/golangci-lint/issues/1779).

Be aware that there may be unused imports after the fix, run `go fmt`.

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
$ testifylint --bool-compare.ignore-custom-types ./...
$ testifylint --expected-actual.pattern=^wanted$ ./...
$ testifylint --formatter.check-format-string --formatter.require-f-funcs ./...
$ testifylint --go-require.ignore-http-handlers ./...
$ testifylint --require-error.fn-pattern="^(Errorf?|NoErrorf?)$" ./...
$ testifylint --suite-extra-assert-call.mode=require ./...
```

### golangci-lint

https://golangci-lint.run/usage/linters/#testifylint

## Checkers

- ✅ – yes
- ❌ – no
- 🤏 – partially

| Name                                                | Enabled By Default | Autofix |
|-----------------------------------------------------|--------------------|---------|
| [blank-import](#blank-import)                       | ✅                  | ❌       |
| [bool-compare](#bool-compare)                       | ✅                  | ✅       |
| [compares](#compares)                               | ✅                  | ✅       |
| [empty](#empty)                                     | ✅                  | ✅       |
| [error-is-as](#error-is-as)                         | ✅                  | 🤏      |
| [error-nil](#error-nil)                             | ✅                  | ✅       |
| [expected-actual](#expected-actual)                 | ✅                  | ✅       |
| [float-compare](#float-compare)                     | ✅                  | ❌       |
| [formatter](#formatter)                             | ✅                  | 🤏      |
| [go-require](#go-require)                           | ✅                  | ❌       |
| [len](#len)                                         | ✅                  | ✅       |
| [negative-positive](#negative-positive)             | ✅                  | ✅       |
| [nil-compare](#nil-compare)                         | ✅                  | ✅       |
| [require-error](#require-error)                     | ✅                  | ❌       |
| [suite-broken-parallel](#suite-broken-parallel)     | ✅                  | ✅       |
| [suite-dont-use-pkg](#suite-dont-use-pkg)           | ✅                  | ✅       |
| [suite-extra-assert-call](#suite-extra-assert-call) | ✅                  | ✅       |
| [suite-subtest-run](#suite-subtest-run)             | ✅                  | ❌       |
| [suite-thelper](#suite-thelper)                     | ❌                  | ✅       |
| [useless-assert](#useless-assert)                   | ✅                  | ❌       |

> ⚠️ Also look at open for contribution [checkers](CONTRIBUTING.md#open-for-contribution)

---

### blank-import

```go
❌
import (
    "testing"

    _ "github.com/stretchr/testify"
    _ "github.com/stretchr/testify/assert"
    _ "github.com/stretchr/testify/http"
    _ "github.com/stretchr/testify/mock"
    _ "github.com/stretchr/testify/require"
    _ "github.com/stretchr/testify/suite"
)

✅
import (
    "testing"
)
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: `testify` doesn't do any `init()` magic, so these imports as `_` do nothing and considered useless.

---

### bool-compare

```go
❌
assert.Equal(t, false, result)
assert.EqualValues(t, false, result)
assert.Exactly(t, false, result)
assert.NotEqual(t, true, result)
assert.NotEqualValues(t, true, result)
assert.False(t, !result)
assert.True(t, result == true)
// And other variations...

✅
assert.True(t, result)
assert.False(t, result)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: Code simplification.

Also `bool-compare` supports user defined types like

```go
type Bool bool
```

And fixes assertions via casting variable to builtin `bool`:

```go
var predicate Bool
❌ assert.Equal(t, false, predicate)
✅ assert.False(t, bool(predicate))
```

To turn off this behavior use the `--bool-compare.ignore-custom-types` flag.

---

### compares

```go
❌
assert.True(t, a == b)
assert.True(t, a != b)
assert.True(t, a > b)
assert.True(t, a >= b)
assert.True(t, a < b)
assert.True(t, a <= b)
assert.False(t, a == b)
// And so on...

✅
assert.Equal(t, a, b)
assert.NotEqual(t, a, b)
assert.Greater(t, a, b)
assert.GreaterOrEqual(t, a, b)
assert.Less(t, a, b)
assert.LessOrEqual(t, a, b)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

If `a` and `b` are pointers then `assert.Same`/`NotSame` is required instead,
due to the inappropriate recursive nature of `assert.Equal` (based on
[reflect.DeepEqual](https://pkg.go.dev/reflect#DeepEqual)).

---

### empty

```go
❌
assert.Len(t, arr, 0)
assert.Equal(t, 0, len(arr))
assert.EqualValues(t, 0, len(arr))
assert.Exactly(t, 0, len(arr))
assert.LessOrEqual(t, len(arr), 0)
assert.GreaterOrEqual(t, 0, len(arr))
assert.Less(t, len(arr), 0)
assert.Greater(t, 0, len(arr))
assert.Less(t, len(arr), 1)
assert.Greater(t, 1, len(arr))
assert.Zero(t, len(arr))
assert.Empty(t, len(arr))

assert.NotEqual(t, 0, len(arr))
assert.NotEqualValues(t, 0, len(arr))
assert.Less(t, 0, len(arr))
assert.Greater(t, len(arr), 0)
assert.Positive(t, len(arr))
assert.NotZero(t, len(arr))
assert.NotEmpty(t, len(arr))

✅
assert.Empty(t, arr)
assert.NotEmpty(t, err)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### error-is-as

```go
❌
assert.Error(t, err, errSentinel) // Typo, errSentinel hits `msgAndArgs`.
assert.NoError(t, err, errSentinel)
assert.True(t, errors.Is(err, errSentinel))
assert.False(t, errors.Is(err, errSentinel))
assert.True(t, errors.As(err, &target))

✅
assert.ErrorIs(t, err, errSentinel)
assert.NotErrorIs(t, err, errSentinel)
assert.ErrorAs(t, err, &target)
```

**Autofix**: partially. <br>
**Enabled by default**: true. <br>
**Reason**: In the first two cases, a common mistake that leads to hiding the incorrect wrapping of sentinel errors.
In the rest cases – more appropriate `testify` API with clearer failure message.

Also `error-is-as` repeats `go vet`'s
[errorsas check](https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/errorsas/errorsas.go)
logic, but without autofix.

---

### error-nil

```go
❌
assert.Nil(t, err)
assert.NotNil(t, err)
assert.Equal(t, nil, err)
assert.EqualValues(t, nil, err)
assert.Exactly(t, nil, err)
assert.ErrorIs(t, err, nil)

assert.NotEqual(t, nil, err)
assert.NotEqualValues(t, nil, err)
assert.NotErrorIs(t, err, nil)

✅
assert.NoError(t, err)
assert.Error(t, err)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### expected-actual

```go
❌
assert.Equal(t, result, expected)
assert.EqualExportedValues(t, resultObj, User{Name: "Rob"})
assert.EqualValues(t, result, 42)
assert.Exactly(t, result, int64(42))
assert.JSONEq(t, result, `{"version": 3}`)
assert.InDelta(t, result, 42.42, 1.0)
assert.InDeltaMapValues(t, result, map[string]float64{"score": 0.99}, 1.0)
assert.InDeltaSlice(t, result, []float64{0.98, 0.99}, 1.0)
assert.InEpsilon(t, result, 42.42, 0.0001)
assert.InEpsilonSlice(t, result, []float64{0.9801, 0.9902}, 0.0001)
assert.IsType(t, result, (*User)(nil))
assert.NotEqual(t, result, "expected")
assert.NotEqualValues(t, result, "expected")
assert.NotSame(t, resultPtr, &value)
assert.Same(t, resultPtr, &value)
assert.WithinDuration(t, resultTime, time.Date(2023, 01, 12, 11, 46, 33, 0, nil), time.Second)
assert.YAMLEq(t, result, "version: '3'")

✅
assert.Equal(t, expected, result)
assert.EqualExportedValues(t, User{Name: "Rob"}, resultObj)
assert.EqualValues(t, 42, result)
// And so on...
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: A common mistake that makes it harder to understand the reason of failed test.

The checker considers the expected value to be a basic literal, constant, or variable whose name matches the pattern
(`--expected-actual.pattern` flag).

It is
planned [to change the order of assertion arguments](https://github.com/stretchr/testify/issues/1089#Argument_order) to
more natural (actual, expected) in `v2` of `testify`.

---

### float-compare

```go
❌
assert.Equal(t, 42.42, result)
assert.EqualValues(t, 42.42, result)
assert.Exactly(t, 42.42, result)
assert.True(t, result == 42.42)
assert.False(t, result != 42.42)

✅
assert.InEpsilon(t, 42.42, result, 0.0001) // Or assert.InDelta
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Do not forget about [floating point rounding issues](https://floating-point-gui.de/errors/comparison/).

This checker is similar to the [floatcompare](https://github.com/golangci/golangci-lint/pull/2608) linter.

---

### formatter

```go
❌
assert.ElementsMatch(t, certConfig.Org, csr.Subject.Org, "organizations not equal")
assert.Error(t, err, fmt.Sprintf("Profile %s should not be valid", test.profile))
assert.Errorf(t, err, fmt.Sprintf("test %s", test.testName))
assert.Truef(t, targetTs.Equal(ts), "the timestamp should be as expected (%s) but was %s", targetTs)
// And other go vet's printf checks...

✅
assert.ElementsMatchf(t, certConfig.Org, csr.Subject.Org, "organizations not equal")
assert.Errorf(t, err, "Profile %s should not be valid", test.profile)
assert.Errorf(t, err, "test %s", test.testName)
assert.Truef(t, targetTs.Equal(ts), "the timestamp should be as expected (%s) but was %s", targetTs, ts)
```

**Autofix**: partially. <br>
**Enabled by default**: true. <br>
**Reason**: Code simplification, protection from bugs, following the "Go Way".

The `formatter` checker has the following features:

#### 1)

Detecting unnecessary `fmt.Sprintf` in assertions. Somewhat reminiscent of the
[staticcheck's S1028 check](https://staticcheck.dev/docs/checks/#S1028).

#### 2)

Validating consistency of assertions format strings and corresponding arguments, using a patched fork of `go vet`'s
[printf](https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/printf/printf.go) analyzer. To
disable this feature, use `--formatter.check-format-string=false` flag.

#### 3)

Requirement of the f-assertions if format string is used. Disabled by default, use `--formatter.require-f-funcs` flag
to enable. This helps follow Go's implicit convention

> Printf-like functions must end with `f`

and sets the stage for moving to `v2` of `testify`. In this way the checker resembles the
[goprintffuncname](https://github.com/jirfag/go-printf-func-name) linter (included in
[golangci-lint](https://golangci-lint.run/usage/linters/)). Also format string in f-assertions is highlighted by IDE
, e.g. GoLand:

<img width="600" alt="F-assertion IDE highlighting" src="https://github.com/Antonboom/testifylint/assets/17127404/9bdab802-d6eb-477d-a411-6cba043d33a5">

#### Historical Reference

<details>

<summary>Click to expand...</summary>

<br>
Those who are new to `testify` may be discouraged by the duplicative API:

```go
func Equal(t TestingT, expected, actual any, msgAndArgs ...any) bool
func Equalf(t TestingT, expected, actual any, msg string, args ...any) bool
```

The f-functions (Equal**f**, Error**f**, etc.) were introduced a long time ago (2017) to close
[uber-go/zap's issue](https://github.com/stretchr/testify/issues/339): `go1.7 vet` complained
about the following logger
[test](https://github.com/uber-go/zap/blame/8f5ee80ab2dbc713823341ce30334cd9c03a98e5/flag_test.go#L60):

```go
if tc.wantErr {
    // flag_test.go:61: possible formatting directive in Error call
    assert.Error(t, err, "Parse(%v) should fail.", tc.args)
    return
}
```

But! It was the incorrect logic inside `go vet`'s **printf** analyzer, and not `testify`'s issue.

Fact is that in Go 1.7 **printf** only
[checked the name of the function](https://github.com/golang/go/blob/2b7a7b710f096b1b7e6f2ab5e9e3ec003ad7cd12/src/cmd/vet/print.go#L69),
but did not take into account its package, thereby reacting to everything that is possible:

```go
// isPrint records the unformatted-print functions.
var isPrint = map[string]bool{
    "error": true,
    "fatal": true,
    // ...
}
```

This behaviour
was [fixed](https://github.com/golang/go/blob/ad7c32dc3b6d5edc3dd72b3e15c80dc4f4c27064/src/cmd/vet/print.go#L62)
in Go 1.10 after a related [issue](https://github.com/golang/go/issues/22936):

```go
// isPrint records the print functions.
var isPrint = map[string]bool{
    "fmt.Errorf": true,
    "fmt.Fprint": true,
    // ...
}
```

Now **printf** only checked Golang standard library functions (unless configured otherwise) and had nothing against
`testify`'s assertions signatures.

Despite this, f-functions have already been released, giving rise to ambiguous API.

But surely the maintainers had no choice to change the signatures in accordance with Go convention, because it would
break backwards compatibility:

```go
func Equal(t TestingT, expected, actual any) bool
func Equalf(t TestingT, expected, actual any, msg string, args ...any) bool
```

But I hope it will be [introduced](https://github.com/stretchr/testify/issues/1089#issuecomment-1695059265)
in `v2` of `testify`.

</details>

---

### go-require

```go
go func() {
    conn, err = lis.Accept()
    require.NoError(t, err) ❌

    if assert.Error(err) {     ✅
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
functions leading to `t.FailNow` (essentially to `runtime.GoExit`) must only be used in the goroutine that runs the
test.
Otherwise, they will not work as declared, namely, finish the test function.

You can disable the `go-require` checker and continue to use `require` as the current goroutine finisher, but this could
lead

1. to possible resource leaks in tests;
2. to increasing of confusion, because functions will be not used as intended.

Typically, any assertions inside goroutines are a marker of poor test architecture.
Try to execute them in the main goroutine and distribute the data necessary for this into it
([example](https://github.com/ipfs/kubo/issues/2043#issuecomment-164136026)).

Also a bad solution would be to simply replace all `require` in goroutines with `assert`
(like
[here](https://github.com/gravitational/teleport/pull/22567/files#diff-9f5fd20913c5fe80c85263153fa9a0b28dbd1407e53da4ab5d09e13d2774c5dbR7377))
– this will only mask the problem.

The checker is enabled by default, because `testinggoroutine` is enabled by default in `go vet`.

In addition, the checker warns about `require` in HTTP handlers (functions and methods whose signature matches
[http.HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc)), because handlers run in a separate
[service goroutine](https://cs.opensource.google/go/go/+/refs/tags/go1.22.3:src/net/http/server.go;l=2782;drc=1d45a7ef560a76318ed59dfdb178cecd58caf948)
that services the HTTP connection. Terminating these goroutines can lead to undefined behaviour and difficulty debugging
tests. You can turn off the check using the `--go-require.ignore-http-handlers` flag.

P.S. Look at [testify's issue](https://github.com/stretchr/testify/issues/772), related to assertions in the goroutines.

---

### len

```go
❌
assert.Equal(t, 3, len(arr))
assert.EqualValues(t, 3, len(arr))
assert.Exactly(t, 3, len(arr))
assert.True(t, len(arr) == 3)

✅
assert.Len(t, arr, 3)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: More appropriate `testify` API with clearer failure message.

---

### negative-positive

```go
❌
assert.Less(t, a, 0)
assert.Greater(t, 0, a)
assert.True(t, a < 0)
assert.True(t, 0 > a)
assert.False(t, a >= 0)
assert.False(t, 0 <= a)

assert.Greater(t, a, 0)
assert.Less(t, 0, a)
assert.True(t, a > 0)
assert.True(t, 0 < a)
assert.False(t, a <= 0)
assert.False(t, 0 >= a)

✅
assert.Negative(t, a)
assert.Positive(t, a)
```

**Autofix**: true. <br>
**Enabled by default**: true <br>
**Reason**: More appropriate `testify` API with clearer failure message.

Typed zeros (like `int8(0)`, ..., `uint64(0)`) are also supported.

---

### nil-compare

```go
❌
assert.Equal(t, nil, value)
assert.EqualValues(t, nil, value)
assert.Exactly(t, nil, value)

assert.NotEqual(t, nil, value)
assert.NotEqualValues(t, nil, value)

✅
assert.Nil(t, value)
assert.NotNil(t, value)
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: Protection from bugs and more appropriate `testify` API with clearer failure message.

Using untyped `nil` in the functions above along with a non-interface type does not make sense:

```go
assert.Equal(t, nil, eventsChan)    // Always fail.
assert.NotEqual(t, nil, eventsChan) // Always pass.
```

The right way:

```go
assert.Equal(t, (chan Event)(nil), eventsChan)
assert.NotEqual(t, (chan Event)(nil), eventsChan)
```

But in the case of `Equal`, `NotEqual` and `Exactly` type casting approach still doesn't work for the function type.

The best option here is to just use `Nil` / `NotNil` (see [details](https://github.com/stretchr/testify/issues/1524)).

---

### require-error

```go
❌
assert.Error(t, err) // s.Error(err), s.Assert().Error(err)
assert.ErrorIs(t, err, io.EOF)
assert.ErrorAs(t, err, &target)
assert.EqualError(t, err, "end of file")
assert.ErrorContains(t, err, "end of file")
assert.NoError(t, err)
assert.NotErrorIs(t, err, io.EOF)

✅
require.Error(t, err) // s.Require().Error(err), s.Require().Error(err)
require.ErrorIs(t, err, io.EOF)
require.ErrorAs(t, err, &target)
// And so on...
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Such "ignoring" of errors leads to further panics, making the test harder to debug.

[testify/require](https://pkg.go.dev/github.com/stretchr/testify@master/require#hdr-Assertions) allows
to stop test execution when a test fails.

By default `require-error` only checks the `*Error*` assertions, presented above. <br>

You can set `--require-error.fn-pattern` flag to limit the checking to certain calls (but still from the list above).
For example, `--require-error.fn-pattern="^(Errorf?|NoErrorf?)$"` will only check `Error`, `Errorf`, `NoError`,
and `NoErrorf`.

Also, to minimize false positives, `require-error` ignores:

- assertions in the `if` condition;
- assertions in the bool expression;
- the entire `if-else[-if]` block, if there is an assertion in any `if` condition;
- the last assertion in the block, if there are no methods/functions calls after it;
- assertions in an explicit goroutine (including `http.Handler`);
- assertions in an explicit testing cleanup function or suite teardown methods;
- sequence of `NoError` assertions.

---

### suite-broken-parallel

```go
func (s *MySuite) SetupTest() {
    s.T().Parallel() ❌
}

// And other hooks...

func (s *MySuite) TestSomething() {
    s.T().Parallel() ❌
    
    for _, tt := range cases {
        s.Run(tt.name, func() {
            s.T().Parallel() ❌
        })
        
        s.T().Run(tt.name, func(t *testing.T) {
            t.Parallel() ❌
        })
    }
}
```

**Autofix**: true. <br>
**Enabled by default**: true. <br>
**Reason**: Protection from undefined behaviour.

`v1` of `testify` doesn't support suite's parallel tests and subtests.

Depending on [case](./analyzer/testdata/src/debug/suite_broken_parallel_test.go) using of `t.Parallel()` leads to

- data race
- panic
- non-working suite hooks
- silent ignoring of this directive

So, `testify`'s maintainers recommend discourage parallel tests inside suite.

<details>
<summary>Related issues...</summary>

- https://github.com/stretchr/testify/issues/187
- https://github.com/stretchr/testify/issues/466
- https://github.com/stretchr/testify/issues/934
- https://github.com/stretchr/testify/issues/1139
- https://github.com/stretchr/testify/issues/1253

</details>

---

### suite-dont-use-pkg

```go
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

### suite-subtest-run

```go
func (s *MySuite) TestSomething() {
    ❌
    s.T().Run("subtest", func(t *testing.T) {
        assert.Equal(t, 42, result)
    })
     
    ✅
    s.Run("subtest", func() {
        s.Equal(42, result)
    }) 
}
```

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Protection from undefined behavior.

According to `testify` [documentation](https://pkg.go.dev/github.com/stretchr/testify/suite#Suite.Run), `s.Run` should
be used for running subtests. This call (among other things) initializes the suite with a fresh instance of `t` and
protects tests from undefined behavior (such as data races).

Autofix is disabled because in the most cases it requires rewriting the assertions in the subtest and can leads to dead
code.

The checker is especially useful in combination with [suite-dont-use-pkg](#suite-dont-use-pkg).

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

The checker rather acts as an example of
a [checkers.AdvancedChecker](https://github.com/Antonboom/testifylint/blob/676324836555445fded4e9afc004101ec6f597fe/internal/checkers/checker.go#L56).

---

### useless-assert

Currently the checker guards against assertion of the same variable:

```go
❌
assert.Equal(t, tt.value, tt.value)
assert.ElementsMatch(t, users, users)
// And so on...
assert.True(t, num > num)
assert.False(t, num == num)
```

More complex cases are [open for contribution](CONTRIBUTING.md#useless-assert).

**Autofix**: false. <br>
**Enabled by default**: true. <br>
**Reason**: Protection from bugs and possible dead code.

---

## Chain of warnings

Linter does not automatically handle the "evolution" of changes.
And in some cases may be dissatisfied with your code several times, for example:

```go
assert.True(err == nil) // compares: use assert.Equal
assert.Equal(t, err, nil) // error-nil: use assert.NoError
assert.NoError(t, err) // require-error: for error assertions use require
require.NoError(t, err)
```

Please [contribute](./CONTRIBUTING.md) if you have ideas on how to make this better.

## testify v2

The second version of `testify` [promises](https://github.com/stretchr/testify/issues/1089) more "pleasant" API and
makes some above checkers irrelevant.

In this case, the possibility of supporting `v2` in the linter is not excluded.

But at the moment it looks like we
are [extremely far](https://github.com/stretchr/testify/issues/1089#issuecomment-1812734472)
from `v2`. Related milestone [here](https://github.com/stretchr/testify/milestone/4).
