package std

import (
	"fmt"

	"github.com/micheledinelli/golamb/common"
)

type Engine struct {
	Config *common.Config
}

func NewEngine(config *common.Config) *Engine {
	return &Engine{Config: config}
}

// EvalSteps fully reduces an expression to its normal form (if it exists)
// using the specified strategy. It can cause an infinite loop if the
// expression does not have a normal form. For example in the untyped
// lambda calculus the expression (\x. x x) (\x. x x) aka omega combinator
// does not have a normal form. Trying to use the evalautes things like
// (\x. x x) (\x. x x) will cause an infinite loop.
//   - If using CallByName or NormalOrder expressions like
//     (\a.\b.a) t ((\x.x x)(\x.x x)) will reduce to t without getting stuck in an infinite loop.
//   - If using CallByValue the same expression will get stuck in an infinite
//     loop trying to reduce the omega combinator.
//
// See [common.Evaluator]
func (e *Engine) EvalSteps(expr common.Expr) (common.Expr, int) {
	ResetFreshCounter()
	steps := 0
	for {
		next, ok := e.step(expr, e.Config.Strategy)
		if !ok {
			return expr, steps
		}
		if e.Config.Trace {
			fmt.Printf("\x1b[32m%d :|>\x1b[0m\x1b[33m %s\x1b[0m\n", steps+1, next.Format())
		}
		expr = next
		steps++
	}
}

// Eval calls [Engine.EvalSteps] and discards the number of steps taken to reduce
// the expression. It returns the fully reduced expression.
func (e *Engine) Eval(expr common.Expr) common.Expr {
	val, _ := e.EvalSteps(expr)
	return val
}
