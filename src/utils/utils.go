package utils

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/micheledinelli/golamb/common"
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
				exprStr = ExpandMacro(exprStr, env)
				env[varName] = exprStr
			}
		}
	}

	return scanner.Err()
}

func ExpandMacro(input string, env map[string]string) string {
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

func ParseArgs() (config *common.Config) {
	config = &common.Config{}

	var strategy string
	var version bool

	flag.StringVar(&strategy, "strategy", "normal", "Evaluation strategy: cbn, cbv, normal, cbpv")
	flag.StringVar(&strategy, "s", "normal", "Evaluation strategy (shorthand)")

	flag.BoolVar(&version, "version", false, "Print version information")
	flag.BoolVar(&version, "v", false, "Print version information (shorthand)")

	flag.BoolVar(&config.BetaSteps, "beta-steps", false, "Show beta-reduction steps")
	flag.BoolVar(&config.BetaSteps, "b", false, "Show beta-reduction steps (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of GoLamb:\n\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")

		flag.PrintDefaults()

		fmt.Fprintln(flag.CommandLine.Output(), "\nExamples:")
		fmt.Fprintln(flag.CommandLine.Output(), "  ./golamb --strategy=cbn")
		fmt.Fprintln(flag.CommandLine.Output(), "  ./golamb -s cbv -b")
	}

	flag.Parse()

	if version {
		fmt.Printf("GoLamb version %s\n", common.Version)
		os.Exit(0)
	}

	switch strings.ToLower(strategy) {
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
		config.CBPVMode = common.CBPVModeCBV
		fmt.Println("reduction strategy: call-by-push-value")
	default:
		fmt.Printf("unknown strategy %q, defaulting to normal order\n", strategy)
		config.Strategy = common.NormalOrder
	}

	return
}
