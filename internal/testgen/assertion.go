package main

type Check struct { // TODO -> Assertion
	Fn    string // "Equal"
	Argsf string // "%s, %s"

	ReportMsgf string //  "use %s.%s"

	ProposedSelector string
	ProposedFn       string // "InDelta"
	ProposedArgsf    string // %s, %s, 0.0001" // TODO -> может не понадобится, как задание со *? точнее не понадобятся argValues в Expand
}
