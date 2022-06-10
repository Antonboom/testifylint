# testifylint
Checks usage of [github.com/stretchr/testify](https://github.com/stretchr/testify).

## Configuring
TODO: via golangci-lint.

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
**Autofix**: yes. <br>
**Enabled by default**: yes.

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
**Autofix**: yes. <br>
**Enabled by default**: yes.

### empty
```go
❌   assert.Len(t, arr, 0)
     assert.Equal(t, len(arr), 0)
     // And other variations around len(arr)...

✅   assert.Empty(t, arr)
     assert.Empty(t, arr)
```
**Autofix**: yes. <br>
**Enabled by default**: yes.

### error
```go
❌   assert.Nil(t, err)
     assert.NotNil(t, err)

✅   assert.NoError(t, err)
     assert.Error(t, err)
```
**Autofix**: yes. <br>
**Enabled by default**: yes.

### error-is
```go
❌   assert.Error(t, err, errSentinel)
     assert.NoError(t, err, errSentinel)

✅   assert.ErrorIs(t, err, errSentinel)
     assert.NotErrorIs(t, err, errSentinel)
```
**Autofix**: yes. <br>
**Enabled by default**: yes.

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
**Autofix**: yes. <br>
**Enabled by default**: yes.

### float-compare
```go
❌   assert.Equal(t, 42.42, a)
     assert.True(t, a == 42.42)
     assert.False(t, a != 42.42)
	
✅   assert.InDelta(t, 42.42, a, 0.0001)
     assert.InEpsilon(t, 42.42, a, 0.01)
     assert.InDelta(t, 42.42, a, 0.001)
```
**Autofix**: no. <br>
**Enabled by default**: yes.

### len
```go
❌   assert.Equal(t, len(arr), 3)
     assert.True(t, len(arr) == 3)

✅   assert.Len(t, arr, 3)
     assert.Len(t, arr, 3)
```
**Autofix**: yes. <br>
**Enabled by default**: yes.

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
**Autofix**: no. <br>
**Enabled by default**: yes.

### suite-dont-use-pkg
```go
❌

✅
```
**Autofix**: -. <br>
**Enabled by default**: -.

### suite-no-extra-assert-call
```go
❌

✅
```
**Autofix**: -. <br>
**Enabled by default**: -.

## Contribution
// - старайтесь, чтобы тестовые файлы не превышали 3к строк
// - реализуйте интерфейс Disabled для выключения
// - сначала пишите генератор тестов
// (я осознаю, что местами тесты избыточны. но считаю, что тестов много не бывает)
// - добавьте тест анализатора
// - потом реализуйте чекер и укажите его в списке
// обновить readme

<details>
  <summary>Open for contribution</summary>

### `empty` checker: zero config
- empty-zero (сделать флаг для строк в линтере)
- 
### `float-compare` checker: structs with float comparisons
поддержка остального дерьма

     assert.NotEqual(t, 42.42, a)
     assert.Greater(t, a, 42.42)
     assert.GreaterOrEqual(t, a, 42.42)
     assert.Less(t, a, 42.42)
     assert.LessOrEqual(t, a, 42.42)
    
assert.True(t, a != 42.42) // assert.False(t, a == 42.42)
assert.True(t, a > 42.42)  // assert.False(t, a <= 42.42)
assert.True(t, a >= 42.42) // assert.False(t, a < 42.42)
assert.True(t, a < 42.42)  // assert.False(t, a >= 42.42)
assert.True(t, a <= 42.42) // assert.False(t, a > 42.42)






### zero

### http-error
http-error (HTTPSuccess + HTTPError)
```go
❌

✅
```
**Autofix**: yes. <br>
**Enabled by default**: no.

### http-code-const
```go
❌

✅
```
**Autofix**: yes. <br>
**Enabled by default**: no.

### negative-positive
```go
❌

✅

❌   assert.Less(t, val, 0)
✅   assert.Negative(t, val)

❌   assert.Greater(t, val, 0)
✅   assert.Positive(t, val)
```
**Autofix**: yes. <br>
**Enabled by default**: yes.

### equal-values
```go
❌   assert.Equal(t, int64(100), price.Amount)
✅   assert.EqualValues(t, 100, price.Amount)
```
**Autofix**: yes. <br>
**Enabled by default**: false.

### no-error-contains
```go
❌   assert.ErrorContains(t, err, "not found")
✅   assert.ErrorIs(t, err, ErrNotFound)
```
**Autofix**: no. <br>
**Enabled by default**: yes.

### no-equal-error
```go
❌   assert.EqualError(t, err, "user not found")
✅   assert.ErrorIs(t, err, ErrUserNotFound)
```
**Autofix**: no. <br>
**Enabled by default**: yes.

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

### suite-test-name
```go
❌

✅
```
**Autofix**: yes. <br>
**Enabled by default**: false.

### suite-thelper
```go
❌

✅
```
**Autofix**: yes. <br>
**Enabled by default**: false.

</details>
