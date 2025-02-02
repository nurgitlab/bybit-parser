package random

import "math/rand"

func NewRandomString(l int) string {
	ans := ""

	for i := 0; i < l; i++ {
		ans += string(byte(randInt(97, 122)))
	}
	return ans
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
