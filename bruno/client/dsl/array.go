package dsl

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type BruArray struct {
	Tag  string
	Data []string
}

func (bt *BruArray) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s [\n", bt.Tag))
	for i, v := range bt.Data {
		suffix := ","
		if i == len(bt.Data)-1 {
			suffix = ""
		}
		retVal.WriteString(fmt.Sprintf("\t%s%s\n", v, suffix))
	}
	retVal.WriteString("]\n")
	return retVal.String()
}

func ImportArray(tag string, scanner *bufio.Scanner, regex map[string]*regexp.Regexp) (*BruArray, error) {
	retVal := BruArray{
		Tag:  tag,
		Data: []string{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		// Skip blank lines
		matched := regex[BLANK].MatchString(line)
		if matched {
			continue
		}
		// Check for block end
		matched = regex[BLOCK_END].MatchString(line)
		if matched {
			return &retVal, nil
		}
		// Check for array item
		matches := regex[ARRAY_ITEM].FindStringSubmatch(line)
		if matches != nil {
			value := matches[1]
			retVal.Data = append(retVal.Data, value)
		}
	}
	return nil, fmt.Errorf("Unexpected end of file while parsing array block: %s", tag)
}
