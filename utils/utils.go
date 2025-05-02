package utils

import (
	"log"
	"math/rand/v2"
	"os"
)

func LoadFile(path string) *os.File {
	file, err := os.Open(path) // Path to your image
	if err != nil {
		log.Fatal(err)
	}

	return file
}

// Utils ---- move to module?
func GetRandomNumInRange(minLimit float64, maxLimit float64) float64 {
	return minLimit + rand.Float64()*(maxLimit-minLimit)
}
