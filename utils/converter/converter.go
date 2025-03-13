package converter

import "strconv"

func StringToInt(input string) int {
	res, _ := strconv.Atoi(input)
	return res
}

func IntToString(input int) string {
	return strconv.Itoa(input)
}
