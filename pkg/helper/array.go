package helper

func IsArrayContainsString(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func IsArrayContainsInt64(arr []int64, number int64) bool {
	for _, item := range arr {
		if item == number {
			return true
		}
	}
	return false
}
