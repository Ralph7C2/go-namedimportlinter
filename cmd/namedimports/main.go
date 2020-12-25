package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ralph7c2/go-namedimportlinter/namedimports"
)

func main() {
	singlechecker.Main(namedimports.NewAnalyzer())
}
