package cbpv

import (
	"fmt"

	"github.com/micheledinelli/golamb/common"
)

type CBPVEngine struct {
	Config *common.Config
}

func NewCBPVEngine(config *common.Config) *CBPVEngine {
	return &CBPVEngine{Config: config}
}

// EvalSteps evaluates an expression by compiling it to Call-by-Push-Value (CBPV)
// intermediate structures and driving a stack-based abstract machine.
// Because CBPV strictly separates Values from Computations, the traditional
// evaluation strategies (Call-by-Value / Call-by-Name) are handled entirely
// by the compilation pass (CBPVMode) rather than changing the abstract machine loop.
//
//   - If compiled under Call-by-Name (lazy), arguments are encapsulated in suspended Thunks
//     and pushed directly onto the stack. Divergent paths like the Omega combinator
//     ((\x. x x) (\x. x x)) are ignored if a function discards its argument, preventing infinite loops.
//
//   - If compiled under Call-by-Value (eager), applications explicitly force arguments to
//     reduce to a native Value before entering the function focus, causing identical divergent
//     expressions to spin into infinite loops.
//
// Instead of recording every micro-transition of the abstract machine (such as stack pushes,
// thunk forcing, or control flow returns), this loop explicitly isolates and counts only
// beta-reductions incrementing the step count only when a computational abstraction (Abs)
// pops and consumes an argument from the activation stack.
//
// See [common.Evaluator]
func (e *CBPVEngine) EvalSteps(expr common.Expr) (common.Expr, int) {
	var stack Stack
	comp := compileToCBPV(expr, e.Config.CBPVMode)
	currentComp := comp

	steps := 0

	for {
		if _, isAbs := currentComp.(*Abs); isAbs && len(stack) > 0 {
			steps++
		}

		next, stepped := step(currentComp, &stack)
		if !stepped {
			switch final := currentComp.(type) {
			case *Return:
				if len(stack) == 0 {
					return toExpr(final.Val), steps
				}

				switch v := final.Val.(type) {
				case *Thunk:
					currentComp = v.Comp
					continue

				case *Var:
					argVal, _ := stack.Pop()
					currentComp = &App{
						Fn:  &Force{Val: v},
						Arg: argVal,
					}
					continue
				}

				panic(fmt.Sprintf("runtime error: stack contains values but computation returned an unhandled value type %T", final.Val))
			case *Abs:
				return toExpr(&Thunk{Comp: final}), steps
			default:
				panic(fmt.Sprintf("runtime error: machine blocked unexpectedly on type %T", final))
			}
		}
		currentComp = next
	}
}

func (e *CBPVEngine) Eval(expr common.Expr) common.Expr {
	val, _ := e.EvalSteps(expr)
	return val
}
