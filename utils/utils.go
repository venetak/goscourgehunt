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
func GetRandomNumInRange(limit float64) float64 {
	return 0 + rand.Float64()*(limit-0)
}
