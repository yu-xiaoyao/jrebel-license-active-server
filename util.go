package main

import (
	"math/rand"
)

const randCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const allRandCharsetLen = len(randCharset)

func newServerRandomness() (serverRandomness string) {
	// H2ulzLlh7E0=
	b := make([]byte, 11)
	for i := 0; i < 11; i++ {
		b[i] = randCharset[rand.Intn(allRandCharsetLen)]
	}
	return string(b) + "="
}
