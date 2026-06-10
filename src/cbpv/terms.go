package cbpv

type Value interface {
	isValue()
	Format() string
}

type Computation interface {
	isComputation()
	Format() string
}

type Var struct{ Name string }
type Thunk struct{ Comp Computation }

func (*Var) isValue()   {}
func (*Thunk) isValue() {}

type Return struct{ Val Value }
type Abs struct {
	Param string
	Body  Computation
}
type App struct {
	Fn  Computation
	Arg Value
}

type Force struct{ Val Value }

func (*Return) isComputation() {}
func (*Abs) isComputation()    {}
func (*App) isComputation()    {}
func (*Force) isComputation()  {}

func (v *Var) Format() string    { return v.Name }
func (t *Thunk) Format() string  { return "thunk(" + t.Comp.Format() + ")" }
func (r *Return) Format() string { return "return " + r.Val.Format() }
func (a *Abs) Format() string    { return "λ" + a.Param + "." + a.Body.Format() }
func (a *App) Format() string    { return "(" + a.Fn.Format() + " " + a.Arg.Format() + ")" }
func (f *Force) Format() string  { return "force(" + f.Val.Format() + ")" }
