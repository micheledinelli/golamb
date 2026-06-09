package lc

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func LoadMacrosFromFile(filePath string, env map[string]string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || (strings.HasPrefix(line, "#")) || (strings.HasPrefix(line, "//")) {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			varName := strings.TrimSpace(parts[0])
			exprStr := strings.TrimSpace(parts[1])

			if varName != "" && !strings.ContainsAny(varName, " \t()\\.") {
				exprStr = ExpandMacroStrings(exprStr, env)
				env[varName] = exprStr
			}
		}
	}

	return scanner.Err()
}

func ExpandMacroStrings(input string, env map[string]string) string {
	for {
		changed := false
		for macroName, macroBody := range env {

			var SB strings.Builder
			runes := []rune(input)

			for i := 0; i < len(runes); {
				if strings.HasPrefix(string(runes[i:]), macroName) {
					endIdx := i + len(macroName)

					startWord := i == 0 || (!unicode.IsLetter(runes[i-1]) && !unicode.IsDigit(runes[i-1]))
					endWord := endIdx == len(runes) || (!unicode.IsLetter(runes[endIdx]) && !unicode.IsDigit(runes[endIdx]))

					if startWord && endWord {
						SB.WriteString("(" + macroBody + ")")
						i = endIdx
						changed = true
						continue
					}
				}
				SB.WriteRune(runes[i])
				i++
			}
			input = SB.String()
		}

		if !changed {
			break
		}
	}
	return input
}

func ParseArgs() (strategy RedStrategy) {
	stratFlag := flag.String("strategy", "normal", "evaluation strategy: cbn, cbv, normal")
	flag.StringVar(stratFlag, "s", "normal", "evaluation strategy (shorthand)")

	versionFlag := flag.Bool("version", false, "print version information")
	flag.BoolVar(versionFlag, "v", false, "print version information (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of GoLamb:\n\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")

		fmt.Fprintln(flag.CommandLine.Output(), "  -s, --strategy string")
		fmt.Fprintln(flag.CommandLine.Output(), "        Evaluation strategy: cbn, cbv, normal (default \"normal\")")
		fmt.Fprintln(flag.CommandLine.Output())

		fmt.Fprintln(flag.CommandLine.Output(), "  -v, --version")
		fmt.Fprintln(flag.CommandLine.Output(), "        Print version information")
		fmt.Fprintln(flag.CommandLine.Output())

		fmt.Fprintln(flag.CommandLine.Output(), "Examples:")
		fmt.Fprintln(flag.CommandLine.Output(), "  ./golamb --strategy=cbn")
		fmt.Fprintln(flag.CommandLine.Output(), "  ./golamb -s cbv")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("GoLamb version %s\n", Version)
		os.Exit(0)
	}

	strategy = NormalOrder
	switch strings.ToLower(*stratFlag) {
	case "cbv", "call-by-value":
		strategy = CallByValue
		fmt.Println("using strategy: Call-by-Value (Eager, strict)")
	case "cbn", "call-by-name":
		strategy = CallByName
		fmt.Println("using strategy: Call-by-Name (Lazy, stops at lambdas)")
	case "normal", "normal-order":
		strategy = NormalOrder
		fmt.Println("using strategy: Normal Order (Lazy, full reduction)")
	default:
		fmt.Printf("unknown strategy %q, defaulting to Normal Order\n", *stratFlag)
		strategy = NormalOrder
	}
	return
}
