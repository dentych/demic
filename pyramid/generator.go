package pyramid

import (
	"bytes"
	"math/rand"
	"time"
)

var validLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	rand.Seed(time.Now().Unix())
}

func GenerateId(amount int) string {
	var buf bytes.Buffer
	for i := 0; i < amount; i++ {
		letter := validLetters[rand.Intn(len(validLetters))]
		buf.WriteByte(letter)
	}
	return buf.String()
}
