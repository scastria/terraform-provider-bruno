package dsl

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type BruDict struct {
	Tag  string
	Data map[string]string
}

func (bd *BruDict) GetTag() string {
	return bd.Tag
}

func (bd *BruDict) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bd.Tag))
	for k, v := range bd.Data {
		retVal.WriteString(fmt.Sprintf("%s%s: %s\n", ITEM_INDENT, k, v))
	}
	retVal.WriteString("}\n")
	return retVal.String()
}

func ImportDict(tag string, scanner *bufio.Scanner, regex map[string]*regexp.Regexp) (*BruDict, error) {
	retVal := BruDict{
		Tag:  tag,
		Data: make(map[string]string),
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
		// Check for dict item
		matches := regex[DICT_ITEM].FindStringSubmatch(line)
		if matches != nil {
			key := matches[1]
			value := matches[2]
			retVal.Data[key] = value
		}
	}
	return nil, fmt.Errorf("Unexpected end of file while parsing dict block: %s", tag)
}
