package main

import (
	"fmt"
	"strings"
)

type Generator interface {
	Generate() map[string]string
	Output(f func(map[string]string))
}
type PhpConstantKeys struct {
	Pattern string
	Keys    []string
}

func (p *PhpConstantKeys) Generate() map[string]string {
	m := make(map[string]string, len(p.Keys))
	for _, v := range p.Keys {
		k := strings.ReplaceAll(v, ".", "_")
		m[fmt.Sprintf(p.Pattern, strings.ToUpper(k))] = v
	}
	return m
}

func (p *PhpConstantKeys) Output(f func(map[string]string)) {
	m := p.Generate()
	if f == nil {
		for k, v := range m {
			fmt.Printf("const %s = \"%s\";\n", k, v)
		}
		return
	}
	f(m)
}
