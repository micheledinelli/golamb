package cbpv

import (
	"github.com/micheledinelli/golamb/common"
)

type CBPVEngine struct {
	Config *common.Config
}

func NewCBPVEngine(config *common.Config) *CBPVEngine {
	return &CBPVEngine{Config: config}
}

// See [common.Evaluator]
func (e *CBPVEngine) EvalSteps(expr common.Expr) (common.Expr, int) {
	// var stack Stack
	// comp := compileToCBPV(expr, e.Config.CBPVMode)
	panic("CBPV evaluation is not yet implemented")
}

// Eval calls [CBPVEngine.EvalSteps] and discards the number of steps taken to reduce
// the expression. It returns the fully reduced expression.
func (e *CBPVEngine) Eval(expr common.Expr) common.Expr {
	val, _ := e.EvalSteps(expr)
	return val
}
