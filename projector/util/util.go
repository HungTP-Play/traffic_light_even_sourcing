package util

import "strconv"

func StringToFloat64(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func StringColorToInt(color string) int {
	switch color {
	case "RED":
		return 1
	case "GREEN":
		return 2
	case "YELLOW":
		return 3
	default:
		return 0
	}
}
