package main

import "fmt"

const (
	sizeB = 1 << (iota * 10)
	sizeKiB
	sizeMiB
	sizeGiB
	sizeTiB
)

func bytesHuman(b int64, precision int) string {
	var v = float64(b)
	var suffix string
	switch {
	case b >= sizeTiB:
		v, suffix = v/sizeTiB, "TiB"
	case b >= sizeGiB:
		v, suffix = v/sizeGiB, "GiB"
	case b >= sizeMiB:
		v, suffix = v/sizeMiB, "MiB"
	case b >= sizeKiB:
		v, suffix = v/sizeKiB, "KiB"
	default:
		v, suffix = v, "B"
	}
	return fmt.Sprintf("%.*f %s", precision, v, suffix)
}
