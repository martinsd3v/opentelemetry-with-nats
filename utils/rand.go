package utils

import (
	"math/rand"
)

func RandNumber(min, max int) int {
	return rand.Intn(max-min) + min
}
