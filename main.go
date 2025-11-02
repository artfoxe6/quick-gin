package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/artfoxe6/quick-gin/internal/scaffold"
)

func main() {
	combinedArgs, err := reorderArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	moduleFlag := fs.String("module", "", "module path for the new project (defaults to project name)")
	forceFlag := fs.Bool("force", false, "overwrite target directory if it already exists")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [flags] <project-name>\n", os.Args[0])
		fs.PrintDefaults()
	}

	if err := fs.Parse(combinedArgs); err != nil {
		if err == flag.ErrHelp {
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	projectName := args[0]
	opts := scaffold.Options{
		ModulePath: *moduleFlag,
		Force:      *forceFlag,
	}

	if err := scaffold.Run(projectName, opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func reorderArgs(raw []string) ([]string, error) {
	var flagArgs []string
	var positional []string

	for i := 0; i < len(raw); i++ {
		arg := raw[i]
		switch {
		case arg == "--module" || arg == "-module":
			if i+1 >= len(raw) {
				return nil, fmt.Errorf("flag %s requires a value", arg)
			}
			flagArgs = append(flagArgs, arg, raw[i+1])
			i++
		case strings.HasPrefix(arg, "--module=") || strings.HasPrefix(arg, "-module="):
			flagArgs = append(flagArgs, arg)
		case arg == "--force" || arg == "-force":
			flagArgs = append(flagArgs, arg)
		case strings.HasPrefix(arg, "--force=") || strings.HasPrefix(arg, "-force="):
			flagArgs = append(flagArgs, arg)
		case arg == "--help" || arg == "-h":
			flagArgs = append(flagArgs, arg)
		default:
			positional = append(positional, arg)
		}
	}

	return append(flagArgs, positional...), nil
}
