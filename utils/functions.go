package utils

// checks if array contains string
func ContainsStr(arr []string, search string) (contains bool) {
	contains = false
	for _, el := range arr {
		if el == search {
			contains = true
			break
		}
	}
	return
}

func RemoveValueFromArray(arr []string, val string) []string {
	var resp []string
	for _, v := range arr {
		if v == val {
			continue
		}
		resp = append(resp, v)
	}
	return resp
}
