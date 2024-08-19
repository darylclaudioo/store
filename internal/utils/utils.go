package utils

func RoundUp(v float64) int {
	if v != float64(int(v)) {
		return int(v) + 1
	}
	return int(v)
}
