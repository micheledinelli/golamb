package cbpv

import (
	"fmt"

	c "github.com/micheledinelli/golamb/common"
)

// step performs a single evaluation step on the given computation,
// using the provided stack for intermediate values. It returns the
// next computation and a boolean indicating whether a step was taken
// (true) or if the computation is in normal form (false).
// It panics if it encounters an invalid state, such as trying to force a
// non-thunk value.
func (e *CBPVEngine) step(comp Computation, stack *Stack) (Computation, bool) {
	switch cv := comp.(type) {
	case *App:
		stack.Push(cv.Arg)
		return cv.Fn, true

	case *Abs:
		argVal, ok := stack.Pop()
		if !ok {
			return cv, false
		}
		if e.Config.CBPVMode == c.CBPVModeCBV {
			if thunk, ok := argVal.(*Thunk); ok {
				argVal = e.forceArg(thunk)
			}
		}
		return substitute(cv.Body, cv.Param, argVal), true

	case *Force:
		if thunk, ok := cv.Val.(*Thunk); ok {
			return thunk.Comp, true
		}
		panic(fmt.Sprintf("runtime error: cannot force non-thunk variant %T", cv.Val))

	case *Return:
		return cv, false
	}
	return comp, false
}

func (e *CBPVEngine) forceArg(thunk *Thunk) Value {
	result, _ := e.EvalSteps(compToExpr(thunk.Comp))
	compiled := compileToCBPV(result, e.Config.CBPVMode)
	if ret, ok := compiled.(*Return); ok {
		return ret.Val
	}
	if abs, ok := compiled.(*Abs); ok {
		return &Thunk{Comp: abs}
	}
	return &Thunk{Comp: compiled}
}

// substitute performs capture-avoiding substitution. It traverses the computation
// structure, replacing occurrences of the target variable with the replacement
// value. It handles variable capture by renaming bound variables when necessary
// to avoid conflicts with free variables in the replacement.
// Freshness is obtained by appending a prime (') to the variable name.
func substitute(comp Computation, target string, replacement Value) Computation {
	switch cv := comp.(type) {
	case *Return:
		return &Return{Val: substValue(cv.Val, target, replacement)}

	case *App:
		return &App{
			Fn:  substitute(cv.Fn, target, replacement),
			Arg: substValue(cv.Arg, target, replacement),
		}

	case *Force:
		return &Force{Val: substValue(cv.Val, target, replacement)}

	case *Abs:
		if cv.Param == target {
			return cv
		}
		if v, ok := replacement.(*Var); ok && cv.Param == v.Name {
			fresh := cv.Param + "'"
			cv.Body = substitute(cv.Body, cv.Param, &Var{Name: fresh})
			cv.Param = fresh
		}
		return &Abs{
			Param: cv.Param,
			Body:  substitute(cv.Body, target, replacement),
		}
	}
	return comp
}

func substValue(val Value, target string, replacement Value) Value {
	switch v := val.(type) {
	case *Var:
		if v.Name == target {
			return replacement
		}
		return v
	case *Thunk:
		return &Thunk{Comp: substitute(v.Comp, target, replacement)}
	}
	return val
}

// toExpr converts a CBPV Value back into a [common.Expr]
// for pretty-printing or further processing. It handles both Var
// and Thunk cases, recursively converting Thunks into their
// underlying computations.
// It panics if it encounters an unknown [Value] type during conversion.
func toExpr(val Value) c.Expr {
	switch v := val.(type) {
	case *Var:
		return &c.Var{Name: v.Name}
	case *Thunk:
		return compToExpr(v.Comp)
	default:
		panic(fmt.Sprintf("unknown cbpv value type during decompilation: %T", val))
	}
}

// compToExpr converts a CBPV Computation back into a [common.Expr]
// by recursively translating each computation form into its
// corresponding expression form.
// It panics if it encounters an unknown [Computation] type during conversion.
func compToExpr(comp Computation) c.Expr {
	switch t := comp.(type) {
	case *Return:
		return toExpr(t.Val)
	case *Force:
		return toExpr(t.Val)
	case *Abs:
		return &c.Abs{
			Param: t.Param,
			Body:  compToExpr(t.Body),
		}
	case *App:
		return &c.App{
			Fn:  compToExpr(t.Fn),
			Arg: toExpr(t.Arg),
		}
	default:
		panic(fmt.Sprintf("unknown cbpv computation type during decompilation: %T", comp))
	}
}
