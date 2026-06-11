package main

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"

	common "github.com/micheledinelli/golamb/common"
	"github.com/micheledinelli/golamb/src/cbpv"
	"github.com/micheledinelli/golamb/src/std"
	utils "github.com/micheledinelli/golamb/src/utils"
)

var (
	runtime *common.Runtime
	config  *common.Config
)

func main() {
	var err error
	config = utils.ParseArgs()

	rl, err := setupScreen()
	if err != nil {
		fmt.Println("error setting up screen:", err)
		return
	}
	defer rl.Close()

	var engine common.Evaluator
	if config.Strategy == common.CallByPushValue {
		engine = cbpv.NewCBPVEngine(config)
		panic("CBPV evaluation is not yet implemented")
	} else {
		engine = std.NewEngine(config)
	}
	runtime = common.NewRuntime(engine)

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Commands
		if line[0] == ':' {
			if err := handleCommand(line, rl); err != nil {
				if err.Error() == "exiting" {
					fmt.Fprintln(rl.Stdout(), "Goodbye!")
					break
				}
				fmt.Fprintln(rl.Stderr(), "command error:", err)
			}
			continue
		}

		if strings.Contains(line, "=") {
			handleAssignment(line, rl)
			continue
		}

		// Parse, resolve and evaluate the expression
		expr, err := std.Parse(line)
		if err != nil {
			fmt.Fprintln(rl.Stderr(), "parse error:", err)
			continue
		}
		expr = runtime.Resolve(expr)
		result, steps := runtime.Engine.EvalSteps(expr)

		fmt.Fprintf(rl.Stdout(), "\x1b[33m%-14s \x1b[0m%s\r\n",
			"normal form:", result.Format(),
		)

		if config.Trace {
			fmt.Fprintf(rl.Stdout(), "\x1b[31m%-14s \x1b[0m%d\r\n",
				"β-reductions:", steps,
			)
		}
	}
}

func handleCommand(line string, rl *readline.Instance) error {
	command := strings.Split(line, " ")[0]
	switch command {
	case ":quit", ":exit", ":q":
		return fmt.Errorf("exiting")
	case ":clear":
		fmt.Fprint(rl.Stdout(), "\033[2J\033[H")
		return nil
	case ":import":
		if after, ok := strings.CutPrefix(line, ":import "); ok {
			filePath := strings.TrimSpace(after)
			if err := utils.Import(runtime, filePath); err != nil {
				fmt.Fprintln(rl.Stderr(), "import error:", err)
			} else {
				fmt.Fprintf(rl.Stdout(), "\x1b[32msuccessfully imported %s\x1b[0m\n", filePath)
			}
		}
		return nil
	case ":env":
		for name, val := range runtime.Env {
			fmt.Fprintf(rl.Stdout(), "\x1b[33m%s = %s\x1b[0m\n", name, val.Format())
		}
		return nil
	case ":reset":
		runtime = common.NewRuntime(runtime.Engine)
		fmt.Fprintln(rl.Stdout(), "environment reset")
		return nil
	case ":info":
		fmt.Fprintf(rl.Stdout(), "\x1b[32mGoLamb version %s\x1b[0m\n", common.Version)
		fmt.Fprintf(rl.Stdout(), "\x1b[32mReduction strategy: %s\x1b[0m\n", config.Strategy)
		return nil
	case ":strat":
		if after, ok := strings.CutPrefix(line, ":strat "); ok {
			strat := strings.TrimSpace(after)
			switch strings.ToLower(strat) {
			case "cbv", "call-by-value":
				config.Strategy = common.CallByValue
				fmt.Println("reduction strategy: call-by-value")
			case "cbn", "call-by-name":
				config.Strategy = common.CallByName
				fmt.Println("reduction strategy: call-by-name")
			case "normal", "normal-order":
				config.Strategy = common.NormalOrder
				fmt.Println("reduction strategy: normal order")
			case "cbpv", "call-by-push-value":
				config.Strategy = common.CallByPushValue
				fmt.Println("reduction strategy: call-by-push-value")
			default:
				fmt.Printf("unknown strategy %q, no changes made\n", strat)
			}
			runtime.Engine = std.NewEngine(config)
		}
		return nil
	}

	return fmt.Errorf("unknown command")
}

func handleAssignment(line string, rl *readline.Instance) {
	parts := strings.SplitN(line, "=", 2)

	varName := strings.TrimSpace(parts[0])
	exprStr := strings.TrimSpace(parts[1])

	if varName == "" || strings.ContainsAny(varName, " \t()\\.") {
		fmt.Fprintln(rl.Stderr(), "invalid assignment: bad variable name")
		return
	}

	expr, err := std.Parse(exprStr)
	if err != nil {
		fmt.Fprintln(rl.Stderr(), "invalid macro definition")
		return
	}

	runtime.Set(varName, expr)

	fmt.Fprintln(rl.Stdout(),
		"\x1b[33m"+varName+" = "+expr.Format()+"\x1b[0m",
	)
}

func setupScreen() (*readline.Instance, error) {
	fmt.Printf(`         __              _    
  __ _ __\ \  __ _ _ __ | |__ 
 / _`+"` "+`/ _ \ \/ _`+"`"+` | '  \| '_ \
 \__, \___/\_\__,_|_|_|_|_.__/
 |___/`+"\x1b[32m v.%s\x1b[0m"+`      
`, common.Version)

	return readline.New("golamb> ")
}
