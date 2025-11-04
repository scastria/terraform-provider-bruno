package dsl

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type BruText struct {
	Tag  string
	Data string
}

func (bt *BruText) GetTag() string {
	return bt.Tag
}

func (bt *BruText) Export() string {
	var retVal strings.Builder
	retVal.WriteString(fmt.Sprintf("%s {\n", bt.Tag))
	for _, l := range strings.Split(bt.Data, "\n") {
		retVal.WriteString(fmt.Sprintf("\t%s\n", l))
	}
	retVal.WriteString("}\n")
	return retVal.String()
}

func ImportText(tag string, scanner *bufio.Scanner, regex map[string]*regexp.Regexp) (*BruText, error) {
	retVal := BruText{
		Tag: tag,
	}
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		// Check for block end
		matched := regex[BLOCK_END].MatchString(line)
		if matched {
			retVal.Data = strings.Join(lines, "\n")
			return &retVal, nil
		}
		lines = append(lines, line)
	}
	return nil, fmt.Errorf("Unexpected end of file while parsing text block: %s", tag)
}
