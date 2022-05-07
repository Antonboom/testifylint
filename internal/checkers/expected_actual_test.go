package checkers

/*
import "testing"

func Test_defaultExpectedVarPattern(t *testing.T) {
	cases := []struct {
		ident   string
		matched bool
	}{
		{ident: "exp", matched: true},
		{ident: "expected", matched: true},
		{ident: "expResult", matched: true},
		{ident: "expectedResult", matched: true},
		{ident: "resultExp", matched: true},
		{ident: "resultExpected", matched: true},

		{ident: "want", matched: true},
		{ident: "wanted", matched: true},
		{ident: "wantError", matched: true},
		{ident: "wantedError", matched: true},
		{ident: "errWant", matched: true},
		{ident: "errWanted", matched: true},

		{ident: "expired", matched: false},
		{ident: "expecting", matched: false},
		{ident: "expresult", matched: false},
		{ident: "expectedresult", matched: false},
		{ident: "resultexp", matched: false},
		{ident: "resultexpected", matched: false},
		{ident: "resultExpires", matched: false},
		{ident: "resultExpectation", matched: false},
		{ident: "wantime", matched: false},
		{ident: "wanteddy", matched: false},
		{ident: "wantresult", matched: false},
		{ident: "wantedresult", matched: false},
		{ident: "resultwant", matched: false},
		{ident: "resultwanted", matched: false},
		{ident: "isClientWantAttention", matched: false},
		{ident: "isClientWantedAttention", matched: false},
		{ident: "clientExpBalance", matched: false},
		{ident: "clientExpectedBalance", matched: false},
	}

	for _, tt := range cases {
		t.Run(tt.ident, func(t *testing.T) {
			if b := defaultExpectedVarPattern.MatchString(tt.ident); b != tt.matched {
				t.Errorf("%q: incorrect regexp", tt.ident)
			}
		})
	}
}
*/
