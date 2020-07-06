package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NamedImports []Import `yaml:"namedImports"`
}

type Import struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

func FromFile(filename string) (Config, error) {
	fBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	var c Config
	err = yaml.Unmarshal(fBytes, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func AddCommandLine(c Config, imports *string) (Config, error) {
	if imports == nil {
		return c, nil
	}
	list := strings.Split(*imports, ",")
	for _, s := range list {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return c, fmt.Errorf("invalid format")
		}
		c.NamedImports = append(c.NamedImports, Import{Path: parts[0], Name: parts[1]})
	}
	return c, nil
}
