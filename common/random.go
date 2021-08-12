package common

import (
	"math/rand"
	"time"
)

const Alphanum string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomString returns a random string of alpha-numeric characters,
// n is the length
func RandomStringAlphanum(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = Alphanum[b%byte(len(Alphanum))]
	}
	return string(bytes)
}

// RandomDuration returns a random duration between 0 and max.
func RandomDuration(max time.Duration) time.Duration {
	if max == 0 {
		return 0
	}

	sleepns := rand.Int63n(max.Nanoseconds())

	return time.Duration(sleepns)
}
