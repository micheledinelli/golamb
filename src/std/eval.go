package std

import (
	"fmt"

	c "github.com/micheledinelli/golamb/common"
)

var freshCounter int

// FV returns the set of free variables in an lambda-expression.
// The set of free variables of a term t, written FV(t), is defined
// inductively on the structure of t as follows:
//   - FV(x) = {x}
//   - FV(M N) = FV(M) U FV(N)
//   - FV(\x.M) = FV(M) \ {x}
func FV(e c.Expr) map[string]struct{} {
	switch t := e.(type) {
	case *c.Var:
		return map[string]struct{}{t.Name: {}}
	case *c.App:
		left := FV(t.Fn)
		right := FV(t.Arg)

		res := make(map[string]struct{})
		for k := range left {
			res[k] = struct{}{}
		}
		for k := range right {
			res[k] = struct{}{}
		}
		return res
	case *c.Abs:
		vars := FV(t.Body)
		res := make(map[string]struct{})
		for k := range vars {
			if k != t.Param {
				res[k] = struct{}{}
			}
		}
		return res
	}
	panic("unreachable")
}

// rename is the alpha-conversion function. It replaces all free occurrences
// of variable old in an expression e with a new fresh variable.
func rename(e c.Expr, old, fresh string) c.Expr {
	switch t := e.(type) {
	case *c.Var:
		if t.Name == old {
			return &c.Var{Name: fresh}
		}
		return t
	case *c.App:
		return &c.App{
			Fn:  rename(t.Fn, old, fresh),
			Arg: rename(t.Arg, old, fresh),
		}
	case *c.Abs:
		param := t.Param
		if param == old {
			param = fresh
		}
		return &c.Abs{
			Param: param,
			Body:  rename(t.Body, old, fresh),
		}
	}
	panic("unreachable")
}

// substitute is the beta-reduction function. It replaces all free occurrences
// of variable x in an expression expr with a value.
// It is defined inductively on the structure of expr as follows:
//   - substitute(x, x, value) = value
//   - substitute(y, x, value) = y if y != x
//   - substitute(M N, x, value) = substitute(M, x, value) substitute(N, x, value)
//   - substitute(\y.M, x, value) = \y.substitute(M, x, value) if y == x
//   - substitute(\y.M, x, value) = \y.substitute(M, x, value) if y != x and y not in FV(value)
//   - substitute(\y.M, x, value) = \z.substitute(substitute(M, y, z), x, value) if y != x and y in FV(value)
func substitute(expr c.Expr, x string, value c.Expr) c.Expr {
	switch t := expr.(type) {
	case *c.Var:
		if t.Name == x {
			return value
		}
		return t
	case *c.App:
		return &c.App{
			Fn:  substitute(t.Fn, x, value),
			Arg: substitute(t.Arg, x, value),
		}
	case *c.Abs:
		if t.Param == x {
			return t
		}
		fv := FV(value)
		if _, ok := fv[t.Param]; ok {
			fresh := freshName(t.Param)
			renamedBody := rename(
				t.Body,
				t.Param,
				fresh,
			)
			return &c.Abs{
				Param: fresh,
				Body: substitute(
					renamedBody,
					x,
					value,
				),
			}
		}
		return &c.Abs{
			Param: t.Param,
			Body: substitute(
				t.Body,
				x,
				value,
			),
		}
	}
	panic("unreachable")
}

// step performs a single reduction step on the given expression according to the specified strategy.
// It returns the new expression and a boolean indicating whether a reduction was performed.
//   - For NormalOrder, it reduces the outermost redex first, allowing reductions inside lambdas.
//   - For CallByName, it reduces the outermost redex first without reducing inside lambdas.
//   - For CallByValue, it first reduces the function and argument to values before applying the function.
func step(expr c.Expr, strat c.Strategy) (c.Expr, bool) {
	switch t := expr.(type) {
	case *c.App:
		if strat == c.CallByValue {
			if next, ok := step(t.Fn, strat); ok {
				return &c.App{Fn: next, Arg: t.Arg}, true
			}
			if next, ok := step(t.Arg, strat); ok {
				return &c.App{Fn: t.Fn, Arg: next}, true
			}
			if abs, ok := t.Fn.(*c.Abs); ok {
				return substitute(abs.Body, abs.Param, t.Arg), true
			}
			return expr, false
		}

		if abs, ok := t.Fn.(*c.Abs); ok {
			return substitute(abs.Body, abs.Param, t.Arg), true
		}
		if next, ok := step(t.Fn, strat); ok {
			return &c.App{Fn: next, Arg: t.Arg}, true
		}
		if next, ok := step(t.Arg, strat); ok {
			return &c.App{Fn: t.Fn, Arg: next}, true
		}
	case *c.Abs:
		if strat == c.CallByName || strat == c.CallByValue {
			return expr, false
		}

		if next, ok := step(t.Body, strat); ok {
			return &c.Abs{Param: t.Param, Body: next}, true
		}
	}
	return expr, false
}

func freshName(base string) string {
	freshCounter++
	return fmt.Sprintf("%s_%d", base, freshCounter)
}

func ResetFreshCounter() {
	freshCounter = 0
}
