# testifylint
Checks usage of [github.com/stretchr/testify](https://github.com/stretchr/testify).

## Installation & usage

```
$ go install github.com/Antonboom/testifylint@latest
$ testifylint ./...
```

## Configuring

```
$ testifylint --enable-all ./...
$ testifylint --enable=empty,error-is ./...
$ testifylint --enable=expected-actual  -expected-actual.pattern=^wanted$ ./...
```

## Checkers

| Name                                                      | Enabled By Default | Autofix |
|-----------------------------------------------------------|--------------------|---------|
| [bool-compare](#bool-compare)                             | ✅                  | ✅       |
| [compares](#compares)                                     | ✅                  | ✅       |
| [empty](#empty)                                           | ✅                  | ✅       |
| [error](#error)                                           | ✅                  | ✅       |
| [error-is](#error-is)                                     | ✅                  | ✅       |
| [expected-actual](#expected-actual)                       | ✅                  | ✅       |
| [float-compare](#float-compare)                           | ✅                  | ❌       |
| [len](#len)                                               | ✅                  | ✅       |
| [require-error](#require-error)                           | ✅                  | ❌       |
| [suite-dont-use-pkg](#suite-dont-use-pkg)                 | ✅                  | ✅       |
| [suite-no-extra-assert-call](#suite-no-extra-assert-call) | ❌                  | ✅       |
| [suite-thelper](#suite-thelper)                           | ❌                  | ✅       |

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
