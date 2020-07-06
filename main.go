package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ralph7c2/namedimportlinter/pkg/config"
	lint "github.com/ralph7c2/namedimportlinter/pkg/linter"
)

func main() {
	conf := flag.String("conf", "", "YAML file with configuration")
	imports := flag.String("imports", "", "Comma separated list of names imports eg. github/ralph7c2/namedimportlinter/pkg/linter:lnt,github.com/ralph7c2/namedimportlinter/config:cfg")
	flag.Parse()

	c := config.Config{}
	if conf != nil {
		c, _ = config.FromFile(*conf)
	}
	c, _ = config.AddCommandLine(c, imports)

	l := lint.NewLinter(&c)

	for _, arg := range flag.Args() {
		l.Lint(arg)
		out, err := l.Output()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if len(out) != 0 {
			fmt.Println(arg)
			fmt.Println(out)
		}
	}
}
