package dsl

import (
	"fmt"
	"strings"
)

type BruText struct {
	Tag  string
	Data string
}

func (bt BruText) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bt.Tag))
	for _, l := range strings.Split(bt.Data, "\n") {
		retVal.WriteString(fmt.Sprintf("\t%s\n", l))
	}
	retVal.WriteString("}\n")
	return retVal.String()
}
