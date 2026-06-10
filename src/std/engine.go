package std

import (
	"fmt"

	c "github.com/micheledinelli/golamb/common"
)

type Engine struct {
	Config *c.Config
}

func NewEngine(config *c.Config) *Engine {
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
// See [c.Evaluator]
func (e *Engine) EvalSteps(expr c.Expr) (c.Expr, int) {
	steps := 0
	for {
		next, ok := step(expr, e.Config.Strategy)
		if !ok {
			return expr, steps
		}
		if e.Config.BetaSteps {
			fmt.Printf("\x1b[33m|>\x1b[0m %v\r\n", expr.Format())
		}
		expr = next
		steps++
	}
}

func (e *Engine) Eval(expr c.Expr) c.Expr {
	val, _ := e.EvalSteps(expr)
	return val
}
