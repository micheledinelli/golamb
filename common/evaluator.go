package common

type Evaluator interface {
	Eval(expr Expr) Expr
	EvalSteps(expr Expr) (Expr, int)
}
