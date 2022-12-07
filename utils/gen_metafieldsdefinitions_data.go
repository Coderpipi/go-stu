package main

import (
	"fmt"
	"strconv"
)

type MetaFieldsDefinitions struct {
	Pattern string
	Keys    [][]any
}

// (159565846631441914, 10234, 'moodyfam 3 handle', 'product', 'my_fields', 'moodyfam_3_handle',
// 'single_line_text_field', '', false, 0)

func (p *MetaFieldsDefinitions) Generate() map[string]string {
	m := make(map[string]string, len(p.Keys))
	for i, key := range p.Keys {
		m[strconv.Itoa(i)] = fmt.Sprintf(p.Pattern, key...)
	}
	return m
}

func (p *MetaFieldsDefinitions) Output(f func(map[string]string)) {

	m := p.Generate()
	if f != nil {
		f(m)
		return
	}
	for _, v := range m {
		fmt.Println(v)
	}
}
