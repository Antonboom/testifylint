package encodedcompareissue196

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoFollowSymlink(t *testing.T) {
	defer PrepareTestEnv(t)()
	session := loginUser(t, "user2")

	assertCase := func(t *testing.T, url, expectedSymlinkURL string, shouldExist bool) {
		t.Helper()

		req := NewRequest(t, "GET", url)
		resp := session.MakeRequest(t, req, http.StatusOK)

		htmlDoc := NewHTMLParser(t, resp.Body)
		symlinkURL, ok := htmlDoc.Find(".file-actions .button[data-kind='follow-symlink']").Attr("href")
		if shouldExist {
			assert.True(t, ok)
			assert.Equal(t, expectedSymlinkURL, symlinkURL)
		} else {
			assert.False(t, ok)
		}
	}

	t.Run("Normal", func(t *testing.T) {
		defer PrintCurrentTest(t)()
		assertCase(t, "/user2/readme-test/src/branch/symlink/README.md?display=source", "/user2/readme-test/src/branch/symlink/some/other/path/awefulcake.txt", true)
	})
}

func PrepareTestEnv(t *testing.T) func()   { return func() {} }
func PrintCurrentTest(t *testing.T) func() { return func() {} }

func loginUser(t *testing.T, uid string) *session {
	return new(session)
}

type session struct {
}

func (s *session) MakeRequest(t *testing.T, req *http.Request, expStatusCode int) *Response {
	return new(Response)
}

func NewRequest(t *testing.T, method string, url string) *http.Request {
	return new(http.Request)
}

type Response struct {
	Body *bytes.Buffer
}

// HTMLDoc struct
type HTMLDoc struct {
}

// NewHTMLParser parse html file
func NewHTMLParser(t testing.TB, body *bytes.Buffer) *HTMLDoc {
	t.Helper()
	return new(HTMLDoc)
}

// Find gets the descendants of each element in the current set of
// matched elements, filtered by a selector. It returns a new Selection
// object containing these matched elements.
func (doc *HTMLDoc) Find(selector string) *Selection {
	return new(Selection)
}

type Selection struct {
}

// Attr gets the specified attribute's value for the first element in the
// Selection. To get the value for each element individually, use a looping
// construct such as Each or Map method.
func (s *Selection) Attr(attrName string) (val string, exists bool) {
	return "", false
}
