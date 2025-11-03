package dsl

import "strings"

type BruDoc struct {
	Data []BruBlock
}

func (bt BruDoc) Export() string {
	var retVal strings.Builder
	for _, block := range bt.Data {
		retVal.WriteString(block.Export())
	}
	return retVal.String()
}
