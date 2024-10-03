package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type EncodedCompareTestsGenerator struct{}

func (EncodedCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewEncodedCompare()
}

func (g EncodedCompareTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	const multiLineJSON = `{
	"id": 123,
	"method": "get_prop",
	"params": ["power","sat"]
}
`
	multiLineCase := "assert.Equal(t, ` // want \"encoded-compare: use assert\\.JSONEq\"\n" + multiLineJSON + "`, raw)"
	multiLineCase += "\nassert.Equal(t, ` // want \"encoded-compare: use assert\\.JSONEq\"\n" + multiLineJSON +
		"`" + `, raw, "msg with args %d %s", 42, "42")`

	const multiLineYAML = `
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
images:
  - name: foo
    newName: bar
  - name: bar
    newName: baz
    newTag: "123"
`

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		MultiLineJSONCase string
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			// Raw strings cases.
			{
				Fn: "Equal", Argsf: "`{\"name\":\"name\",\"value\":1000}`, respBody",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "expBody, respBody",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "`{\"status\":404,\"message\":\"abc\"}`, string(respBytes)",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "`{\"message\":\"success\"}`, w.Body.String()",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: `"{\n  \"first\": \"Tobi\",\n  \"last\": \"Ferret\"\n}", string(w.Body.Bytes())`,
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: `"{\n  \"first\": \"Tobi\",\n  \"last\": \"Ferret\"\n}", w.Body.String()`,
			},
			{
				Fn: "Equal", Argsf: `"{\n\t\"msg\": \"hello world\"\n}", respBody`,
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "fmt.Sprintf(`{\"value\":\"%s\",\"valuePtr\":\"%s\"}`, hexString, hexString), string(respBytes)",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},

			// Variable name cases.
			{
				Fn: "Equal", Argsf: "`[{\"@id\":\"a\",\"b\":[{\"@id\":\"c\"}]}]`, toJSON",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: `"{\"FirstName\":\"john\",\"LastName\":\"doe\",\"Age\":26,\"Height\":182.88}", string(resJson)`,
				ProposedArgsf: `"{\"FirstName\":\"john\",\"LastName\":\"doe\",\"Age\":26,\"Height\":182.88}", resJson`,
				ReportMsgf:    report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "expJSON, resultJSON",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "jsonb, respBody",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "respBody, jsonb",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "expJSON, resJson",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},

			{
				Fn: "Equal", Argsf: "expectedYAML, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "expYaml, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "ymlResult, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "yamlResult, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "expYML, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "conf, expectedYAML",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Equal", Argsf: "outputYaml, string(output.Bytes())",
				ReportMsgf: report, ProposedFn: "YAMLEq",
				ProposedArgsf: "outputYaml, output.String()",
			},

			// Type conversion cases.
			{
				Fn: "Equal", Argsf: "json.RawMessage(`{\"uuid\": \"b65b1a22-db6d-4f5a-9b3d-7302368a82e6\"}`), batch.ParentSummary()",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "`{\"uuid\": \"b65b1a22-db6d-4f5a-9b3d-7302368a82e6\"}`, string(batch.ParentSummary())",
			},
			{
				Fn: "Equal", Argsf: "res[0].Data, json.RawMessage([]byte(`{\"name\":\"new\"}`))",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "string(res[0].Data), `{\"name\":\"new\"}`",
			},
			{
				Fn: "Equal", Argsf: "json.RawMessage(raw), json.RawMessage(resultJSONBytes)",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "raw, string(resultJSONBytes)",
			},
			{
				Fn: "Equal", Argsf: "json.RawMessage(raw), raw",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "raw, raw",
			},
			{
				Fn: "Equal", Argsf: `json.RawMessage("{}"), respBody`,
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: `"{}", respBody`,
			},
			{
				Fn: "Equal", Argsf: `respBody, json.RawMessage("null")`,
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: `respBody, "null"`,
			},
			{
				Fn: "Equal", Argsf: "json.RawMessage(`[\"more\",\"raw\",\"things\"]`), resultJSONBytes",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "`[\"more\",\"raw\",\"things\"]`, string(resultJSONBytes)",
			},

			{
				Fn: "Equal", Argsf: `"{}", string(resultJSONBytes)`,
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Equal", Argsf: "[]byte(expJSON), resultJSONBytes",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "expJSON, string(resultJSONBytes)",
			},
			{
				Fn: "Equal", Argsf: "[]byte(expYaml), respBytes",
				ReportMsgf: report, ProposedFn: "YAMLEq",
				ProposedArgsf: "expYaml, string(respBytes)",
			},

			// Replace/Trim cases.
			{
				Fn: "Equal", Argsf: `expJSON, strings.Trim(string(resultJSONBytes), "\n")`,
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "expJSON, string(resultJSONBytes)",
			},
			{
				Fn: "Equal", Argsf: `raw, strings.Replace(jsonb, "\n", "", -1)`,
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "raw, jsonb",
			},
			{
				Fn: "Equal", Argsf: "`{\"status\":\"healthy\",\"message\":\"\",\"peer_count\":1}`," +
					" strings.ReplaceAll(string(respBytes), \" \", \"\")",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "`{\"status\":\"healthy\",\"message\":\"\",\"peer_count\":1}`, string(respBytes)",
			},
			{
				Fn: "Equal", Argsf: "`{\"foo\":\"bar\"}`, strings.TrimSpace(w.Body.String())",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "`{\"foo\":\"bar\"}`, w.Body.String()",
			},
			{
				Fn: "Equal", Argsf: "`{\"bar\":\"foo\"}`, strings.TrimSpace(string(w.Body.Bytes()))",
				ReportMsgf: report, ProposedFn: "JSONEq",
				ProposedArgsf: "`{\"bar\":\"foo\"}`, w.Body.String()",
			},
			{
				Fn: "Equal", Argsf: `strings.TrimSpace(strings.ReplaceAll(expYaml, "\t", "  ")), strings.TrimSpace(string(respBytes))`,
				ReportMsgf: report, ProposedFn: "YAMLEq",
				ProposedArgsf: "expYaml, string(respBytes)",
			},

			// Other Equal* cases.
			{
				Fn: "EqualValues", Argsf: "expJSON, resJson",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "Exactly", Argsf: "expJSON, resJson",
				ReportMsgf: report, ProposedFn: "JSONEq",
			},
			{
				Fn: "EqualValues", Argsf: "expYaml, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
			{
				Fn: "Exactly", Argsf: "expYaml, conf",
				ReportMsgf: report, ProposedFn: "YAMLEq",
			},
		},
		MultiLineJSONCase: multiLineCase,
		ValidAssertions: []Assertion{
			{Fn: "JSONEq", Argsf: "`{\"name\":\"name\",\"value\":1000}`, respBody"},
			{Fn: "JSONEq", Argsf: "expJSON, resultJSON"},
			{Fn: "JSONEq", Argsf: "`{\"foo\":\"bar\"}`, `{\"foo\":\"bar\"}`"},
			{Fn: "JSONEq", Argsf: "`{\"message\":\"success\"}`, w.Body.String()"},
			{Fn: "JSONEq", Argsf: "fmt.Sprintf(`{\"value\":\"%s\"}`, hexString), resJson"},
			{Fn: "JSONEq", Argsf: `"{\n  \"first\": \"Tobi\",\n  \"last\": \"Ferret\"\n}", w.Body.String()`},

			{Fn: "YAMLEq", Argsf: "expYaml, conf"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: `"{{ .StepName }}", "use", "command name incorrect"`},
			{Fn: "Equal", Argsf: "json.RawMessage{}, respBody"},
			{Fn: "Equal", Argsf: "json.RawMessage(nil), respBody"},

			{Fn: "Equal", Argsf: "raw, raw"},
			{Fn: "EqualValues", Argsf: "raw, raw"},
			{Fn: "Exactly", Argsf: "raw, raw"},
			{Fn: "JSONEq", Argsf: "raw, raw"},

			{Fn: "Equal", Argsf: "string(respBytes), raw"},
			{Fn: "EqualValues", Argsf: "raw, string(respBytes)"},
			{Fn: "Exactly", Argsf: "string(respBytes), raw"},
			{Fn: "JSONEq", Argsf: "raw, string(respBytes)"},

			{Fn: "NotEqual", Argsf: "raw, resultJSON"},
			{Fn: "NotEqualValues", Argsf: "resultJSON, resultJSON"},

			{Fn: "YAMLEq", Argsf: "`" + multiLineYAML + "`, conf"},                // Not supported.
			{Fn: "YAMLEq", Argsf: `"kind: Kustomization", "kind: Kustomization"`}, // Not supported.
			{Fn: "YAMLEq", Argsf: "raw, conf"},
			{Fn: "YAMLEq", Argsf: "raw, string(respBytes)"},
		},
	}
}

func (EncodedCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("EncodedCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(encodedCompareTestTmpl))
}

func (EncodedCompareTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("EncodedCompareTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(encodedCompareTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const encodedCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var respBody, raw, hexString, toJSON, expJSON, resultJSON, jsonb, resJson string
	var conf, expectedYAML, expYaml, ymlResult, yamlResult, expYML, outputYaml string
	var respBytes, resultJSONBytes []byte
	w := httptest.NewRecorder()
	var batch interface { ParentSummary() []byte }
	var res [1]struct{ Data []byte }
	var output bytes.Buffer

	const expBody = ` + "`{\"status\":\"healthy\",\"message\":\"\",\"peer_count\":1}`" + `

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
		{{ .MultiLineJSONCase }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Ignored.
	{
		{{- range $ai, $assrn := $.IgnoredAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}
`
