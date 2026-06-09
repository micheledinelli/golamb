package lc

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
func (a *Abs) Format() string { return "\\" + a.Param + "." + a.Body.Format() }
func (a *App) Format() string { return "(" + a.Fn.Format() + " " + a.Arg.Format() + ")" }
