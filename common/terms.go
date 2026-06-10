package common

import "fmt"

type Expr interface {
	expr()
	Format() string
}

type Var struct {
	Name string
}

type Abs struct {
	Param string
	Body  Expr
}

type App struct {
	Fn  Expr
	Arg Expr
}

func (*Var) expr() {}
func (*Abs) expr() {}
func (*App) expr() {}

func (v *Var) Format() string { return v.Name }
func (a *Abs) Format() string { return fmt.Sprintf("λ%s. %s", a.Param, a.Body.Format()) }
func (a *App) Format() string {
	var fnStr, argStr string

	if _, isAbs := a.Fn.(*Abs); isAbs {
		fnStr = fmt.Sprintf("(%s)", a.Fn.Format())
	} else {
		fnStr = a.Fn.Format()
	}

	_, isApp := a.Arg.(*App)
	_, isAbs := a.Arg.(*Abs)
	if isApp || isAbs {
		argStr = fmt.Sprintf("(%s)", a.Arg.Format())
	} else {
		argStr = a.Arg.Format()
	}

	return fmt.Sprintf("%s %s", fnStr, argStr)
}
