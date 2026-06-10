package tests

import (
	"testing"
	"time"

	"github.com/micheledinelli/golamb/common"
	"github.com/micheledinelli/golamb/src/std"
	"github.com/micheledinelli/golamb/src/utils"
)

func TestRuntimeAssignment(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.NormalOrder})
	runtime := common.NewRuntime(engine)

	// Store \x.x in the runtime under the name "id"
	runtime.Set("id", mustParse(`\x.x`))

	// Now we can resolve "id" to its definition and evaluate it
	expr := mustParse("id")
	expr = runtime.Resolve(expr)

	// Evaluating "id" should give us back the identity
	result := runtime.Engine.Eval(expr)

	if result.Format() != "λx. x" {
		t.Fatalf("expected λx. x, got %s", result.Format())
	}
}

func TestBooleanLogic(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.NormalOrder})
	runtime := common.NewRuntime(engine)

	runtime.Set("t", mustParse(`\t.\f.t`))
	runtime.Set("f", mustParse(`\t.\f.f`))

	runtime.Set("and", mustParse(`\p.\q.p q p`))
	runtime.Set("or", mustParse(`\p.\q.p p q`))
	runtime.Set("not", mustParse(`\p.p f t`))
	runtime.Set("xor", mustParse(`\p.\q.p (not q) q`))

	andExpr := runtime.Resolve(mustParse("and t f"))
	andResult := runtime.Engine.Eval(andExpr)
	if andResult.Format() != "λt. λf. f" {
		t.Fatalf("expected false, got %s", andResult.Format())
	}

	orExpr := runtime.Resolve(mustParse("or t f"))
	orResult := runtime.Engine.Eval(orExpr)
	if orResult.Format() != "λt. λf. t" {
		t.Fatalf("expected true, got %s", orResult.Format())
	}

	notExpr := runtime.Resolve(mustParse("not t"))
	notResult := runtime.Engine.Eval(notExpr)
	if notResult.Format() != "λt. λf. f" {
		t.Fatalf("expected false, got %s", notResult.Format())
	}

	xorExpr := runtime.Resolve(mustParse("xor t f"))
	xorResult := runtime.Engine.Eval(xorExpr)
	if xorResult.Format() != "λt. λf. t" {
		t.Fatalf("expected true, got %s", xorResult.Format())
	}
}

func TestImport(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.NormalOrder})
	runtime := common.NewRuntime(engine)

	utils.Import(runtime, "std.test")

	expr := mustParse("ID")
	expr = runtime.Resolve(expr)

	result := runtime.Engine.Eval(expr)
	if result.Format() != "λx. x" {
		t.Fatalf("expected λx. x, got %s", result.Format())
	}
}

func TestChurchNumerals(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.NormalOrder})
	runtime := common.NewRuntime(engine)

	utils.Import(runtime, "std.test")

	addExpr := runtime.Resolve(mustParse("ADD 1 2"))
	addResult := runtime.Engine.Eval(addExpr)
	if addResult.Format() != "λf. λx. f (f (f x))" {
		t.Fatalf("expected 3, got %s", addResult.Format())
	}

	mulExpr := runtime.Resolve(mustParse("MULT 2 2"))
	mulResult := runtime.Engine.Eval(mulExpr)
	if mulResult.Format() != "λf. λx. f (f (f (f x)))" {
		t.Fatalf("expected 4, got %s", mulResult.Format())
	}

	succExpr := runtime.Resolve(mustParse("SUCC 1"))
	succResult := runtime.Engine.Eval(succExpr)
	if succResult.Format() != "λf. λx. f (f x)" {
		t.Fatalf("expected 2, got %s", succResult.Format())
	}

	predExpr := runtime.Resolve(mustParse("PRED 2"))
	predResult := runtime.Engine.Eval(predExpr)
	if predResult.Format() != "λf. λx. f x" {
		t.Fatalf("expected 1, got %s", predResult.Format())
	}

	isZeroExpr := runtime.Resolve(mustParse("ISZERO 0"))
	isZeroResult := runtime.Engine.Eval(isZeroExpr)
	if isZeroResult.Format() != "λt. λf. t" {
		t.Fatalf("expected true, got %s", isZeroResult.Format())
	}
}

func TestCBVStuck(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.CallByValue})
	runtime := common.NewRuntime(engine)

	utils.Import(runtime, "std.test")

	expr := mustParse(`(\x.\y.x) w ((\x.x x)(\x.x x))`)
	expr = runtime.Resolve(expr)

	done := make(chan struct{})

	go func() {
		runtime.Engine.Eval(expr)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("expected divergence, but evaluation terminated")
	case <-time.After(2 * time.Second):
		// diverges (or at least did not finish quickly)
	}
}

func TestCBNTerminates(t *testing.T) {
	engine := std.NewEngine(&common.Config{Strategy: common.CallByName})
	runtime := common.NewRuntime(engine)

	utils.Import(runtime, "std.test")

	expr := mustParse(`(\x.\y.x) w ((\x.x x)(\x.x x))`)
	expr = runtime.Resolve(expr)

	result := runtime.Engine.Eval(expr)
	if result.Format() != "w" {
		t.Fatalf("expected w, got %s", result.Format())
	}
}

func mustParse(s string) common.Expr {
	expr, err := std.Parse(s)
	if err != nil {
		panic("failed to parse expression: " + s)
	}
	return expr
}
