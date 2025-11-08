package dsl

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type BruMultiDict struct {
	Tag  string
	Data map[string][]string
}

func (bd *BruMultiDict) GetTag() string {
	return bd.Tag
}

func (bd *BruMultiDict) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bd.Tag))
	for k, v := range bd.Data {
		for _, vv := range v {
			retVal.WriteString(fmt.Sprintf("%s%s: %s\n", ITEM_INDENT, k, vv))
		}
	}
	retVal.WriteString("}\n")
	return retVal.String()
}

func ImportMultiDict(tag string, scanner *bufio.Scanner, regex map[string]*regexp.Regexp) (*BruMultiDict, error) {
	retVal := BruMultiDict{
		Tag:  tag,
		Data: make(map[string][]string),
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
			curValue, exists := retVal.Data[key]
			if exists {
				curValue = append(curValue, value)
				retVal.Data[key] = curValue
			} else {
				retVal.Data[key] = []string{value}
			}
		}
	}
	return nil, fmt.Errorf("Unexpected end of file while parsing multi dict block: %s", tag)
}
