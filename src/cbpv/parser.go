package cbpv

import (
	c "github.com/micheledinelli/golamb/common"
)

// compileToCBPV transforms a lambda expression into its corresponding CBPV computation form
// based on the specified CBPVMode (Call-by-Value or Call-by-Name).
func compileToCBPV(expr c.Expr, mode c.CBPVMode) Computation {
	if mode == c.CBPVModeCBV {
		return compileCBV(expr)
	}
	return compileCBN(expr)
}

// Compiles a lambda expression into CBPV form under Call-by-Value semantics,
// where arguments are fully reduced to Values before being passed to functions.
// It panics if it encounters an unknown expression type during compilation.
func compileCBV(expr c.Expr) Computation {
	switch e := expr.(type) {
	case *c.Var:
		return &Return{Val: &Var{Name: e.Name}}

	case *c.Abs:
		return &Return{
			Val: &Thunk{
				Comp: &Abs{
					Param: e.Param,
					Body:  compileCBV(e.Body),
				},
			},
		}
	case *c.App:
		return &App{
			Fn:  &Force{Val: &Thunk{Comp: compileCBV(e.Fn)}},
			Arg: &Thunk{Comp: compileCBV(e.Arg)},
		}
	}
	panic("unknown expression type")
}

// Compiles a lambda expression into CBPV form under Call-by-Name (lazy)
// semantics, where arguments are suspended in Thunks and forced only when consumed by an Abs.
// It panics if it encounters an unknown expression type during compilation.
func compileCBN(expr c.Expr) Computation {
	switch e := expr.(type) {
	case *c.Var:
		return &Force{Val: &Var{Name: e.Name}}

	case *c.Abs:
		return &Abs{
			Param: e.Param,
			Body:  compileCBN(e.Body),
		}

	case *c.App:
		return &App{
			Fn:  &Force{Val: &Thunk{Comp: compileCBV(e.Fn)}},
			Arg: &Thunk{Comp: compileCBV(e.Arg)},
		}
	}
	panic("unknown expression type")
}
