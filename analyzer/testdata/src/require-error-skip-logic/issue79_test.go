package requireerrorskiplogic

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertRegexpError1(regexp any) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, got error, msg ...any) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		return assert.Error(t, got, msg...) && assert.Regexp(t, regexp, got.Error(), msg...)
	}
}

func assertRegexpError2(regexp any) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, got error, msg ...any) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		assert.Error(t, got, msg...) // want "require-error: for error assertions use require"
		return assert.Regexp(t, regexp, got.Error(), msg...)
	}
}

func assertRegexpError3(regexp any) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, got error, msg ...any) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		return assert.Error(t, got, msg...) &&
			(assert.Regexp(t, regexp, got.Error(), msg...) ||
				assert.ErrorContains(t, got, "failed to"))
	}
}

func requireRegexpError1(regexp any) require.ErrorAssertionFunc {
	return func(t require.TestingT, got error, msg ...any) {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		assert.Error(t, got, msg...) // want "require-error: for error assertions use require"
		assert.Regexp(t, regexp, got.Error(), msg...)
	}
}

func requireRegexpError2(regexp any) require.ErrorAssertionFunc {
	return func(t require.TestingT, got error, msg ...any) {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		require.Error(t, got, msg...)
		assert.Regexp(t, regexp, got.Error(), msg...)
	}
}
