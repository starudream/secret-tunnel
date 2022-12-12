package unitx

import (
	"fmt"
)

var m = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

func getSizeAndUnit(size float64, base float64, _map []string) (float64, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}

func HumanSizeWithPrecision(size float64, precision int) string {
	size, unit := getSizeAndUnit(size, 1000.0, m)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

func HumanSize(size float64) string {
	return HumanSizeWithPrecision(size, 4)
}
