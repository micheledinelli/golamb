package common

type Strategy int

const (
	NormalOrder     Strategy = iota // Outermost first, reduces inside lambdas
	CallByName                      // Outermost first, stops at lambdas (no reducing bodies)
	CallByValue                     // Evaluates arguments to a value BEFORE applying the function
	CallByPushValue                 // Uses the CBPV stack machine for evaluation
)

type CBPVMode int

const (
	CBPVModeCBN CBPVMode = iota // Compile to CBPV in a call-by-name style
	CBPVModeCBV                 // Compile to CBPV in a call-by-value style
)
