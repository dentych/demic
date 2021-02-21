package util

func ArrayContains(array []string, contain string) bool {
	for _, v := range array {
		if v == contain {
			return true
		}
	}
	return false
}