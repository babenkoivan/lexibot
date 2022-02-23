package utils

import (
	"math/rand"
	"time"
)

func SourcedRand() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}
