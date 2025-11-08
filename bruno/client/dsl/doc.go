package dsl

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	DICT_TAG        = "dict"
	MULTI_DICT_TAG  = "multiDict"
	TEXT_TAG        = "text"
	ARRAY_TAG       = "array"
	DISABLED_PREFIX = "~"
	ITEM_INDENT     = "  "
	BLANK           = `^\s*$`
	BLOCK_START     = `^(?P<tag>\S+)\s*[{\[]$`
	BLOCK_END       = `^[}\]]$`
	DICT_ITEM       = `^\s*(?P<key>\S+)\s*:\s*(?P<value>.*)\s*$`
	ARRAY_ITEM      = `^\s*(?P<value>.+?),?\s*$`
)

type BruDoc struct {
	Data []BruBlock
}

func (bd *BruDoc) GetBlock(tag string) (BruBlock, error) {
	for _, block := range bd.Data {
		if block.GetTag() == tag {
			return block, nil
		}
	}
	return nil, fmt.Errorf("Block with tag %s not found", tag)
}

func (bd *BruDoc) ExportDoc(filePath string) error {
	err := os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	var retVal strings.Builder
	for _, block := range bd.Data {
		retVal.WriteString(block.Export())
	}
	err = os.WriteFile(filePath, []byte(retVal.String()), 0644)
	return err
}

func ImportDoc(filePath string, schema map[string]string) (*BruDoc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	blocks := []BruBlock{}
	regex := map[string]*regexp.Regexp{
		BLANK:       regexp.MustCompile(BLANK),
		BLOCK_START: regexp.MustCompile(BLOCK_START),
		BLOCK_END:   regexp.MustCompile(BLOCK_END),
		DICT_ITEM:   regexp.MustCompile(DICT_ITEM),
		ARRAY_ITEM:  regexp.MustCompile(ARRAY_ITEM),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip blank lines
		matched := regex[BLANK].MatchString(line)
		if matched {
			continue
		}
		// Check for block start
		matches := regex[BLOCK_START].FindStringSubmatch(line)
		if matches != nil {
			tag := matches[1]
			switch schema[tag] {
			case DICT_TAG:
				dictBlock, err := ImportDict(tag, scanner, regex)
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, dictBlock)
				break
			case MULTI_DICT_TAG:
				multiDictBlock, err := ImportMultiDict(tag, scanner, regex)
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, multiDictBlock)
				break
			case TEXT_TAG:
				textBlock, err := ImportText(tag, scanner, regex)
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, textBlock)
				break
			case ARRAY_TAG:
				arrayBlock, err := ImportArray(tag, scanner, regex)
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, arrayBlock)
				break
			default:
				return nil, fmt.Errorf("Unknown block tag: %s", tag)
			}
		}
	}
	retVal := &BruDoc{
		Data: blocks,
	}
	return retVal, nil
}
