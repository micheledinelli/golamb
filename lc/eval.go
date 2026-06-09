package lc

import "fmt"

type RedStrategy int

var freshCounter int

const (
	NormalOrder RedStrategy = iota // Outermost first, reduces inside lambdas
	CallByName                     // Outermost first, stops at lambdas (no reducing bodies)
	CallByValue                    // Evaluates arguments to a value BEFORE applying the function
)

// FV returns the set of free variables in an lambda-expression.
// The set of free variables of a term t, written FV(t), is defined
// inductively on the structure of t as follows:
//   - FV(x) = {x}
//   - FV(M N) = FV(M) U FV(N)
//   - FV(\x.M) = FV(M) \ {x}
func FV(e Expr) map[string]struct{} {
	switch t := e.(type) {
	case *Var:
		return map[string]struct{}{t.Name: {}}
	case *App:
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
	case *Abs:
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

// Rename is the alpha-conversion function. It replaces all free occurrences
// of variable old in an expression e with a new fresh variable.
func Rename(e Expr, old, fresh string) Expr {
	switch t := e.(type) {
	case *Var:
		if t.Name == old {
			return &Var{Name: fresh}
		}
		return t
	case *App:
		return &App{
			Fn:  Rename(t.Fn, old, fresh),
			Arg: Rename(t.Arg, old, fresh),
		}
	case *Abs:
		param := t.Param
		if param == old {
			param = fresh
		}
		return &Abs{
			Param: param,
			Body:  Rename(t.Body, old, fresh),
		}
	}
	panic("unreachable")
}

// Substitute is the beta-reduction function. It replaces all free occurrences
// of variable x in an expression expr with a value.
// It is defined inductively on the structure of expr as follows:
//   - Substitute(x, x, value) = value
//   - Substitute(y, x, value) = y if y != x
//   - Substitute(M N, x, value) = Substitute(M, x, value) Substitute(N, x, value)
//   - Substitute(\y.M, x, value) = \y.Substitute(M, x, value) if y == x
//   - Substitute(\y.M, x, value) = \y.Substitute(M, x, value) if y != x and y not in FV(value)
//   - Substitute(\y.M, x, value) = \z.Substitute(Substitute(M, y, z), x, value) if y != x and y in FV(value)
func Substitute(expr Expr, x string, value Expr) Expr {
	switch t := expr.(type) {
	case *Var:
		if t.Name == x {
			return value
		}
		return t
	case *App:
		return &App{
			Fn:  Substitute(t.Fn, x, value),
			Arg: Substitute(t.Arg, x, value),
		}
	case *Abs:
		if t.Param == x {
			return t
		}
		fv := FV(value)
		if _, ok := fv[t.Param]; ok {
			fresh := FreshName(t.Param)
			renamedBody := Rename(
				t.Body,
				t.Param,
				fresh,
			)
			return &Abs{
				Param: fresh,
				Body: Substitute(
					renamedBody,
					x,
					value,
				),
			}
		}
		return &Abs{
			Param: t.Param,
			Body: Substitute(
				t.Body,
				x,
				value,
			),
		}
	}
	panic("unreachable")
}

// Step performs a single reduction step on the given expression according to the specified strategy.
// It returns the new expression and a boolean indicating whether a reduction was performed.
//   - For NormalOrder, it reduces the outermost redex first, allowing reductions inside lambdas.
//   - For CallByName, it reduces the outermost redex first without reducing inside lambdas.
//   - For CallByValue, it first reduces the function and argument to values before applying the function.
func Step(expr Expr, strat RedStrategy) (Expr, bool) {
	switch t := expr.(type) {
	case *App:
		if strat == CallByValue {
			if next, ok := Step(t.Fn, strat); ok {
				return &App{Fn: next, Arg: t.Arg}, true
			}
			if next, ok := Step(t.Arg, strat); ok {
				return &App{Fn: t.Fn, Arg: next}, true
			}
			if abs, ok := t.Fn.(*Abs); ok {
				return Substitute(abs.Body, abs.Param, t.Arg), true
			}
			return expr, false
		}

		if abs, ok := t.Fn.(*Abs); ok {
			return Substitute(abs.Body, abs.Param, t.Arg), true
		}
		if next, ok := Step(t.Fn, strat); ok {
			return &App{Fn: next, Arg: t.Arg}, true
		}
		if next, ok := Step(t.Arg, strat); ok {
			return &App{Fn: t.Fn, Arg: next}, true
		}
	case *Abs:
		if strat == CallByName || strat == CallByValue {
			return expr, false
		}

		if next, ok := Step(t.Body, strat); ok {
			return &Abs{Param: t.Param, Body: next}, true
		}
	}
	return expr, false
}

// Normalize fully reduces an expression to its normal form (if it exists)
// using the specified strategy. It can cause an infinite loop if the
// expression does not have a normal form. For example in the untyped
// lambda calculus the expression (\x. x x) (\x. x x) aka omega combinator
// does not have a normal form. Trying to use the evalautes things like
// (\x. x x) (\x. x x) will cause an infinite loop.
//   - If using CallByName or NormalOrder expressions like
//     (\a.\b.a) t ((\x.x x)(\x.x x)) will reduce to t without getting stuck in an infinite loop.
//   - If using CallByValue the same expression will get stuck in an infinite
//     loop trying to reduce the omega combinator.
func Normalize(expr Expr, strat RedStrategy) Expr {
	for {
		next, ok := Step(expr, strat)
		if !ok {
			return expr
		}
		expr = next
	}
}

func FreshName(base string) string {
	freshCounter++
	return fmt.Sprintf("%s_%d", base, freshCounter)
}

func ResetFreshCounter() {
	freshCounter = 0
}
