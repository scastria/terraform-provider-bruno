package dsl

import (
	"fmt"
	"strings"
)

type BruDict struct {
	Tag  string
	Data map[string]interface{}
}

func (bt BruDict) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bt.Tag))
	for k, v := range bt.Data {
		retVal.WriteString(fmt.Sprintf("\t%s: %v\n", k, v))
	}
	retVal.WriteString("}\n")
	return retVal.String()
}
