package utils

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
