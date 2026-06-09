package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	lc "github.com/micheledinelli/golamb/lc"
	"golang.org/x/term"
)

func main() {
	var err error
	var oldState *term.State
	var strategy lc.RedStrategy = lc.ParseArgs()
	var screen *term.Terminal = term.NewTerminal(os.Stdin, "> ")

	// Global macro environment for assignments and imports
	var env map[string]string = map[string]string{}

	if oldState, err = term.MakeRaw(int(os.Stdin.Fd())); err != nil {
		fmt.Println("error setting raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		var line string
		var expr lc.Expr

		if line, err = screen.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, "read error:", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "exit" || line == "quit" {
			break
		}
		if line == "" {
			continue
		}

		lc.ResetFreshCounter()

		// Check for :import command
		if after, ok := strings.CutPrefix(line, ":import "); ok {
			filePath := strings.TrimSpace(after)
			if err := lc.LoadMacrosFromFile(filePath, env); err != nil {
				fmt.Fprintf(screen, "import error: %v\r\n", err)
			}

			fmt.Fprintf(screen, "successfully imported macros from %s\r\n", filePath)
			continue
		}

		// Check for assignments
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			varName := strings.TrimSpace(parts[0])
			exprStr := strings.TrimSpace(parts[1])

			if varName == "" || strings.ContainsAny(varName, " \t()\\.") {
				fmt.Fprintln(screen, "invalid assignment: variable name cannot be empty or contain spaces or special characters")
				continue
			}

			exprStr = lc.ExpandMacroStrings(exprStr, env)
			env[varName] = exprStr
			fmt.Fprintf(screen, "defined %s\r\n", varName)
			continue
		}

		expandedLine := lc.ExpandMacroStrings(line, env)
		if expr, err = lc.Parse(expandedLine); err != nil {
			fmt.Fprintln(screen, "parse error:", err)
			continue
		}

		result := lc.Normalize(expr, strategy)
		cleanResult := lc.Normalize(result, lc.NormalOrder)

		fmt.Fprintln(screen, cleanResult.Format())
	}
}
