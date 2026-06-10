package utils

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/micheledinelli/golamb/common"
	"github.com/micheledinelli/golamb/src/std"
)

func ParseArgs() (config *common.Config) {
	config = &common.Config{}

	var strategy string
	var version bool

	flag.StringVar(&strategy, "strategy", "normal", "Evaluation strategy: cbn, cbv, normal, cbpv")
	flag.StringVar(&strategy, "s", "normal", "Evaluation strategy (shorthand)")

	flag.BoolVar(&version, "version", false, "Print version information")
	flag.BoolVar(&version, "v", false, "Print version information (shorthand)")

	flag.BoolVar(&config.Trace, "beta-steps", false, "Show beta-reduction steps")
	flag.BoolVar(&config.Trace, "b", false, "Show beta-reduction steps (shorthand)")

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
	case "cbn", "call-by-name":
		config.Strategy = common.CallByName
	case "normal", "normal-order":
		config.Strategy = common.NormalOrder
	case "cbpv", "call-by-push-value":
		config.Strategy = common.CallByPushValue
		config.CBPVMode = common.CBPVModeCBV
	default:
		fmt.Printf("unknown strategy %q, defaulting to normal order\n", strategy)
		config.Strategy = common.NormalOrder
	}

	return
}

// Import reads a file containing macro definitions and adds them to the provided runtime.
// The file should have lines in the format: name = expression, where name is the macro name
// and expression is a lambda expression parsable by [std.Parse].
// It ignores empty lines and lines starting with # or // (comments).
// It returns an error if the file cannot be read.
func Import(r *common.Runtime, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" ||
			strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "//") {
			continue
		}

		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		name := strings.TrimSpace(parts[0])
		exprStr := strings.TrimSpace(parts[1])

		if name == "" || strings.ContainsAny(name, " \t()\\.") {
			return fmt.Errorf("invalid macro name: %s", name)
		}

		expr, err := std.Parse(exprStr)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		r.Set(name, expr)
	}

	return scanner.Err()
}
