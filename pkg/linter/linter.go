package linter

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ralph7c2/go-namedimportlinter/pkg/config"
)

type linter struct {
	fset         *token.FileSet
	out          io.ReadWriter
	namedImports map[string]string
}

func NewLinter(cfg *config.Config) *linter {
	m := make(map[string]string)
	for _, namedImport := range cfg.NamedImports {
		m[namedImport.Path] = namedImport.Name
	}
	return &linter{
		fset:         token.NewFileSet(),
		out:          bytes.NewBuffer([]byte{}),
		namedImports: m,
	}
}

func (l *linter) Lint(fileName string) {
	f, err := parser.ParseFile(l.fset, fileName, nil, parser.DeclarationErrors)
	if err != nil {
		panic(fmt.Sprintf("Parser error: %s", err))
	}
	for _, spec := range f.Imports {
		path := strings.Trim(spec.Path.Value, "\"")
		name, ok := l.namedImports[path]
		if !ok {
			continue
		}
		if spec.Name == nil {
			l.NoName(path, name)
			continue
		}
		if spec.Name.Name != name {
			l.WrongName(path, name, spec.Name.Name)
			continue
		}
	}
}

func (l *linter) NoName(path, name string) {
	fmt.Fprintf(l.out, "Expected name [%s] for import [%s], but is unnamed\n", name, path)
}

func (l *linter) WrongName(path, expectedName, name string) {
	fmt.Fprintf(l.out, "Expected name [%s] for import [%s], but found name [%s]\n", expectedName, path, name)
}

func (l linter) Output() (string, error) {
	buf, err := ioutil.ReadAll(l.out)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
