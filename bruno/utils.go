package bruno

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
