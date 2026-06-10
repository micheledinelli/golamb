package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	common "github.com/micheledinelli/golamb/common"
	"github.com/micheledinelli/golamb/src/cbpv"
	"github.com/micheledinelli/golamb/src/std"
	utils "github.com/micheledinelli/golamb/src/utils"
	"golang.org/x/term"
)

var (
	env map[string]string
)

func main() {
	var err error
	var oldState *term.State
	var config *common.Config = utils.ParseArgs()
	var screen *term.Terminal = term.NewTerminal(os.Stdin, "golamb> ")

	PrintWelcome(screen)

	// Global macro environment for assignments and imports
	env = map[string]string{}

	if oldState, err = term.MakeRaw(int(os.Stdin.Fd())); err != nil {
		fmt.Println("error setting raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var engine common.Evaluator
	if config.Strategy == common.CallByPushValue {
		engine = cbpv.NewCBPVEngine(config)
	} else {
		engine = std.NewEngine(config)
	}

	for {
		var line string
		var expr common.Expr

		if line, err = screen.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, "read error:", err)
			break
		}

		if line = strings.TrimSpace(line); line == "" {
			continue
		}

		// Handle commands
		if line[0] == ':' {
			if err = handleCommand(line, screen); err != nil {
				if err.Error() == "exiting" {
					fmt.Fprintln(screen, "Goodbye!")
					break
				}
				fmt.Fprintln(screen, "command error:", err)
			}
			continue
		}

		std.ResetFreshCounter()

		// Handle assignments
		if strings.Contains(line, "=") {
			handleAssignment(line, screen)
			continue
		}

		input := utils.ExpandMacro(line, env)
		expr, err = std.Parse(input)

		if err != nil {
			fmt.Fprintln(screen, "parse error:", err)
			continue
		}

		result, steps := engine.EvalSteps(expr)

		fmt.Fprintf(screen, "\x1b[33m%-14s \x1b[0m%s\r\n", "normal form:", result.Format())
		if config.BetaSteps {
			fmt.Fprintf(screen, "\x1b[31m%-14s \x1b[0m%d\r\n", "β-reductions:", steps)
		}
	}
}

func handleCommand(line string, screen *term.Terminal) error {
	command := strings.Split(line, " ")[0]
	switch command {
	case ":quit", ":exit", ":q":
		return fmt.Errorf("exiting")
	case ":import":
		if after, ok := strings.CutPrefix(line, ":import "); ok {
			filePath := strings.TrimSpace(after)
			if err := utils.LoadMacrosFromFile(filePath, env); err != nil {
				fmt.Fprintf(screen, "import error: %v\r\n", err)
				return nil
			}
			fmt.Fprintf(screen, "\x1b[33msuccessfully imported macros from %s\x1b[0m\r\n", filePath)
			return nil
		}
	case ":env":
		for key, value := range env {
			fmt.Fprintf(screen, "\x1b[33m%s = %s\x1b[0m\r\n", key, value)
		}
		return nil
	default:
		return fmt.Errorf("unknown command")
	}
	return nil
}

func handleAssignment(line string, screen *term.Terminal) {
	parts := strings.SplitN(line, "=", 2)
	varName := strings.TrimSpace(parts[0])
	exprStr := strings.TrimSpace(parts[1])

	if varName == "" || strings.ContainsAny(varName, " \t()\\.") {
		fmt.Fprintln(screen, "invalid assignment: variable name cannot be empty or contain spaces or special characters")
		return
	}

	exprStr = utils.ExpandMacro(exprStr, env)
	_, err := std.Parse(exprStr)
	if err != nil {
		fmt.Fprintln(screen, "invalid macro definition: ", err)
		return
	}

	env[varName] = exprStr
	fmt.Fprintln(screen, "\x1b[33m"+varName+" = "+exprStr+"\x1b[0m")
}

func PrintWelcome(screen *term.Terminal) {
	fmt.Fprintf(screen, `
             _              _     
  ____  ___ | | _____ ____ | |__  
 / _  |/ _ \| |(____ |    \|  _ \ 
( (_| | |_| | |/ ___ | | | | |_) )
 \___ |\___/ \_)_____|_|_|_|____/ 
(_____|`+"\x1b[32m Version %s\x1b[0m"+`
`, common.Version)
	fmt.Println()

}
