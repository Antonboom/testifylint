package t

//func TestConfusedWithExpectedActual(t *testing.T) {
//	//type user struct{ name string }
//	//
//	//var (
//	//	i            int
//	//	s            string
//	//	expected     *user
//	//	expectedUser *user
//	//)
//	//
//	//assert.Les(t)
//
//	// f-functions
//}
//
//func TestComparisonInsteadOfEqual(t *testing.T) {
//
//}
//
//func TestComparisonInsteadOfGreater(t *testing.T) {
//
//}
//
//func TestComparisonInsteadOfGreaterOrEqual(t *testing.T) {
//
//}
//
//func TestComparisonInsteadOfLess(t *testing.T) {
//
//}
//
//func TestComparisonInsteadOfLessOrEqual(t *testing.T) {
//
//}
//
//func TestNotNilInsteadOfError(t *testing.T) {
//	// f-functions
//}
//
//func TestNilInsteadOfNoError(t *testing.T) {
//	err := operation()
//
//	assert.Nil(t, err)                         // want "for a better message use assert.NoError instead"
//	assert.Nilf(t, err, "msg")                 // want "for a better message use assert.NoErrorf instead"
//	assert.Nilf(t, err, "msg with arg %d", 42) // want "for a better message use assert.NoErrorf instead"
//
//	require.Nil(t, err)                         // want "for a better message use require.NoError instead"
//	require.Nilf(t, err, "msg")                 // want "for a better message use require.NoErrorf instead"
//	require.Nilf(t, err, "msg with arg %d", 42) // want "for a better message use require.NoErrorf instead"
//
//	assert.NoError(t, err)
//	assert.NoErrorf(t, err, "msg")
//	assert.NoErrorf(t, err, "msg with arg %d", 42)
//
//	require.NoError(t, err)
//	require.NoErrorf(t, err, "msg")
//	require.NoErrorf(t, err, "msg with arg %d", 42)
//}
//
//func TestErrorInsteadOfErrorIs(t *testing.T) {
//	err := operation()
//
//	assert.Error(t, err, io.EOF)                                    // want "invalid usage of assert.Error, use assert.ErrorIs instead"
//	assert.Error(t, err, new(os.PathError))                         // want "invalid usage of assert.Error, use assert.ErrorIs instead"
//	assert.Error(t, err, errors.New("sky is falling"))              // want "invalid usage of assert.Error, use assert.ErrorIs instead"
//	assert.Error(t, err, fmt.Errorf("sky is falling %d times", 10)) // want "invalid usage of assert.Error, use assert.ErrorIs instead"
//
//	require.Error(t, err, io.EOF)                                    // want "invalid usage of require.Error, use require.ErrorIs instead"
//	require.Error(t, err, new(os.PathError))                         // want "invalid usage of require.Error, use require.ErrorIs instead"
//	require.Error(t, err, errors.New("sky is falling"))              // want "invalid usage of require.Error, use require.ErrorIs instead"
//	require.Error(t, err, fmt.Errorf("sky is falling %d times", 10)) // want "invalid usage of require.Error, use require.ErrorIs instead"
//
//	// f-functions
//}
//
//func TestFloatCompare(t *testing.T) {
//	// assert.Equal(t, 1. == 2.)
//
//	// f-functions
//}
//
//func TestBoolAsserts(t *testing.T) {
//
//	// f-functions
//}

// у обычных функций тоже есть msgAndArgs ((

//func TestEmptyAsserts(t *testing.T) {
//	// f-functions
//}
//
//func TestHTTPStatusCode(t *testing.T) {
//	// f-functions
//}
//
//func operation() error {
//	return errors.New("invalid")
//}
//
//// t.helper
//// "don't use Equal for float" or "use NoError instead of Nil for error check"
