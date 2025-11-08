package bruno

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bruno/bruno/client/dsl"
)

func convertInterfaceArrayToStringArray(arr []interface{}) []string {
	retVal := []string{}
	for _, a := range arr {
		line := ""
		if a != nil {
			line = a.(string)
		}
		retVal = append(retVal, line)
	}
	return retVal
}

func createVariableDictBlockFromMap(tag string, variables *schema.Set) *dsl.BruDict {
	retVal := dsl.BruDict{
		Tag:  tag,
		Data: make(map[string]string),
	}
	for _, variable := range variables.List() {
		variableMap := variable.(map[string]interface{})
		varPrefix := ""
		if variableMap["disabled"].(bool) {
			varPrefix = dsl.DISABLED_PREFIX
		}
		retVal.Data[fmt.Sprintf("%s%s", varPrefix, variableMap["key"].(string))] = variableMap["value"].(string)
	}
	return &retVal
}

func createQueryParamMultiDictBlockFromMap(tag string, variables *schema.Set) *dsl.BruMultiDict {
	retVal := dsl.BruMultiDict{
		Tag:  tag,
		Data: make(map[string][]string),
	}
	for _, variable := range variables.List() {
		variableMap := variable.(map[string]interface{})
		varPrefix := ""
		if variableMap["disabled"].(bool) {
			varPrefix = dsl.DISABLED_PREFIX
		}
		fullKey := fmt.Sprintf("%s%s", varPrefix, variableMap["key"].(string))
		curValue, exists := retVal.Data[fullKey]
		if exists {
			curValue = append(curValue, variableMap["value"].(string))
			retVal.Data[fullKey] = curValue
		} else {
			retVal.Data[fullKey] = []string{variableMap["value"].(string)}
		}
	}
	return &retVal
}

//	func createQueryParamMultiDictBlockFromMap(tag string, variables *schema.Set) *dsl.BruMultiDict {
//		retVal := dsl.BruMultiDict{
//			Tag:  tag,
//			Data: make(map[string][]string),
//		}
//		for _, variableSet := range variables.List() {
//			variableSetMap := variableSet.(map[string]interface{})
//			key := variableSetMap["key"].(string)
//			variableSetSet := variableSetMap["values"].(*schema.Set)
//			for _, variable := range variableSetSet.List() {
//				variableMap := variable.(map[string]interface{})
//				varPrefix := ""
//				if variableMap["disabled"].(bool) {
//					varPrefix = dsl.DISABLED_PREFIX
//				}
//				fullKey := fmt.Sprintf("%s%s", varPrefix, key)
//				curValue, exists := retVal.Data[fullKey]
//				if exists {
//					curValue = append(curValue, variableMap["value"].(string))
//					retVal.Data[fullKey] = curValue
//				} else {
//					retVal.Data[fullKey] = []string{variableMap["value"].(string)}
//				}
//			}
//		}
//		return &retVal
//	}
func createTextBlockFromArray(tag string, arr []interface{}) *dsl.BruText {
	retVal := dsl.BruText{
		Tag: tag,
	}
	lines := convertInterfaceArrayToStringArray(arr)
	retVal.Data = strings.Join(lines, "\n")
	return &retVal
}

func createMapFromVariableDictBlock(variableDict *dsl.BruDict) []map[string]interface{} {
	var retVal []map[string]interface{}
	for k, v := range variableDict.Data {
		variableMap := make(map[string]interface{})
		variableMap["key"] = strings.TrimPrefix(k, dsl.DISABLED_PREFIX)
		variableMap["value"] = v
		variableMap["disabled"] = strings.HasPrefix(k, dsl.DISABLED_PREFIX)
		retVal = append(retVal, variableMap)
	}
	return retVal
}

func createMapFromQueryParamMultiDictBlock(variableDict *dsl.BruMultiDict) []map[string]interface{} {
	var retVal []map[string]interface{}
	for k, v := range variableDict.Data {
		for _, vv := range v {
			variableMap := make(map[string]interface{})
			variableMap["key"] = strings.TrimPrefix(k, dsl.DISABLED_PREFIX)
			variableMap["value"] = vv
			variableMap["disabled"] = strings.HasPrefix(k, dsl.DISABLED_PREFIX)
			retVal = append(retVal, variableMap)
		}
	}
	return retVal
}

func createArrayFromTextBlock(textBlock *dsl.BruText) []string {
	return strings.Split(textBlock.Data, "\n")
}

func prepareFolderName(name string) string {
	// Remove any leading slashes
	name = strings.TrimLeft(name, "/")
	// Replace any remaining slashes with hyphens
	name = strings.ReplaceAll(name, "/", "-")
	return name
}
