package dsl

import (
	"fmt"
	"strings"
)

type BruArray struct {
	Tag  string
	Data []string
}

func (bt BruArray) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bt.Tag))
	for _, v := range bt.Data {
		retVal.WriteString(fmt.Sprintf("\t%s\n", v))
	}
	retVal.WriteString("}\n")
	return retVal.String()
}
