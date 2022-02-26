package char

import (
	"bytes"
	"strings"
)

type MappingCharFilter struct {
	old [][]byte
	new [][]byte
}

func NewMappingCharFilter(mappings []string) *MappingCharFilter {
	m := &MappingCharFilter{}
	for _, field := range mappings {
		field := strings.Split(field, " => ")
		if len(field) != 2 {
			continue
		}
		m.old = append(m.old, []byte(field[0]))
		m.new = append(m.new, []byte(field[1]))
	}

	return m
}

func (t *MappingCharFilter) Filter(input []byte) []byte {
	for i := 0; i < len(t.old); i++ {
		input = []byte(bytes.ReplaceAll(input, t.old[i], t.new[i]))
	}
	return input
}
