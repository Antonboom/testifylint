// Code generated by testifylint/internal/testgen. DO NOT EDIT.

package blankimport

import "testing"
import _ "github.com/stretchr/testify"         // want "blank-import: avoid blank import of github.com/stretchr/testify as it does nothing"
import _ "github.com/stretchr/testify/assert"  // want "blank-import: avoid blank import of github.com/stretchr/testify/assert as it does nothing"
import _ "github.com/stretchr/testify/http"    // want "blank-import: avoid blank import of github.com/stretchr/testify/http as it does nothing"
import _ "github.com/stretchr/testify/mock"    // want "blank-import: avoid blank import of github.com/stretchr/testify/mock as it does nothing"
import _ "github.com/stretchr/testify/require" // want "blank-import: avoid blank import of github.com/stretchr/testify/require as it does nothing"
import _ "github.com/stretchr/testify/suite"   // want "blank-import: avoid blank import of github.com/stretchr/testify/suite as it does nothing"
import "strings"

import (
	_ "github.com/stretchr/testify"         // want "blank-import: avoid blank import of github.com/stretchr/testify as it does nothing"
	_ "github.com/stretchr/testify/assert"  // want "blank-import: avoid blank import of github.com/stretchr/testify/assert as it does nothing"
	_ "github.com/stretchr/testify/http"    // want "blank-import: avoid blank import of github.com/stretchr/testify/http as it does nothing"
	_ "github.com/stretchr/testify/mock"    // want "blank-import: avoid blank import of github.com/stretchr/testify/mock as it does nothing"
	_ "github.com/stretchr/testify/require" // want "blank-import: avoid blank import of github.com/stretchr/testify/require as it does nothing"
	_ "github.com/stretchr/testify/suite"   // want "blank-import: avoid blank import of github.com/stretchr/testify/suite as it does nothing"
)

import (
	_ "gopkg.in/yaml.v3"
	"net/url"

	// Testing.
	_ "github.com/stretchr/testify"         // want "blank-import: avoid blank import of github.com/stretchr/testify as it does nothing"
	_ "github.com/stretchr/testify/assert"  // want "blank-import: avoid blank import of github.com/stretchr/testify/assert as it does nothing"
	_ "github.com/stretchr/testify/http"    // want "blank-import: avoid blank import of github.com/stretchr/testify/http as it does nothing"
	_ "github.com/stretchr/testify/mock"    // want "blank-import: avoid blank import of github.com/stretchr/testify/mock as it does nothing"
	_ "github.com/stretchr/testify/require" // want "blank-import: avoid blank import of github.com/stretchr/testify/require as it does nothing"
	_ "github.com/stretchr/testify/suite"   // want "blank-import: avoid blank import of github.com/stretchr/testify/suite as it does nothing"

	. "database/sql"
	_ "github.com/pmezard/go-difflib/difflib"
)

import (
	_ "github.com/stretchr/testify" // want "blank-import: avoid blank import of github.com/stretchr/testify as it does nothing"
)

import (
	// Test dependencies so that it doesn't get cleaned by glide vc
	_ "github.com/stretchr/testify" // want "blank-import: avoid blank import of github.com/stretchr/testify as it does nothing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDummy(t *testing.T) {
	dummy := 1 + 3
	if dummy != 4 {
		t.Errorf("expected %d, but got %d", 4, dummy)
	}

	_ = strings.Builder{}
	_ = url.URL{}
	_ = DB{}

	_ = assert.Equal
	_ = http.TestRoundTripper{}
	_ = mock.Mock{}
	_ = require.Equal
	_ = suite.Suite{}
}
