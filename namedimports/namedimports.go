package namedimports

import (
	"flag"
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/ralph7c2/go-namedimportlinter/pkg/config"
)

func NewAnalyzer() *analysis.Analyzer {
	flagset := flag.NewFlagSet("flags", flag.ExitOnError)
	conf := flagset.String("conf", "", "YAML file with configuration")
	imports := flagset.String("imports", "", "Comma separated list of names imports eg. github/ralph7c2/namedimportlinter/pkg/linter:lnt,github.com/ralph7c2/namedimportlinter/config:cfg")

	flag.Parse()

	c := config.Config{}
	if conf != nil {
		c, _ = config.FromFile(*conf)
	}

	c, _ = config.AddCommandLine(c, imports)

	m := make(map[string]string)
	for _, namedImport := range c.NamedImports {
		m[namedImport.Path] = namedImport.Name
	}

	return &analysis.Analyzer{
		Name: "namedimports",
		Doc:  "ensures defined named imports are consistent",
		Run:  run(m),
		Flags: *flagset,
	}
}

func run(namedImports map[string]string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		for s, s2 := range namedImports {
			fmt.Println(s, s2)
		}

		for _, file := range pass.Files {
			for _, spec := range file.Imports {
				path := strings.Trim(spec.Path.Value, "\"")

				name, ok := namedImports[path]
				if !ok {
					continue
				}

				if spec.Name == nil {
					pass.Reportf(spec.Pos(), "Expected name [%s] for import [%s], but is unnamed", name, path)
					continue
				}

				if spec.Name.Name != name {
					pass.Reportf(spec.Pos(), "Expected name [%s] for import [%s], but found name [%s]\n", name, path, spec.Name.Name)
					continue
				}
			}
		}

		return nil, nil
	}
}
