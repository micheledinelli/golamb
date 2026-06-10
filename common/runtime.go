package common

import "maps"

type Runtime struct {
	Env    map[string]Expr
	Engine Evaluator
}

func NewRuntime(engine Evaluator) *Runtime {
	return &Runtime{
		Env:    make(map[string]Expr),
		Engine: engine,
	}
}

// Resolve takes an expression and recursively resolves all variables
// in it using the runtime's environment.
func (r *Runtime) Resolve(e Expr) Expr {
	return r.resolve(e, map[string]bool{})
}

func (r *Runtime) resolve(e Expr, bound map[string]bool) Expr {
	switch t := e.(type) {
	case *Var:
		if bound[t.Name] {
			return t
		}
		if def, ok := r.Env[t.Name]; ok {
			return r.resolve(def, bound)
		}
		return t

	case *App:
		return &App{
			Fn:  r.resolve(t.Fn, bound),
			Arg: r.resolve(t.Arg, bound),
		}

	case *Abs:
		newBound := copyBound(bound)
		newBound[t.Param] = true

		return &Abs{
			Param: t.Param,
			Body:  r.resolve(t.Body, newBound),
		}
	}
	return e
}

func (r *Runtime) Set(name string, expr Expr) {
	r.Env[name] = expr
}

func copyBound(b map[string]bool) map[string]bool {
	n := make(map[string]bool)
	maps.Copy(n, b)
	return n
}
