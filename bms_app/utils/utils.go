package utils

import (
	"crypto/rand"
	"syscall/js"
	"time"
)

func GenerateRandomId(this js.Value, args []js.Value) interface{} {
	id, err := rand.Prime(rand.Reader, 20)
	if err != nil {
		id, err = rand.Prime(rand.Reader, 20)
	}

	return int(id.Uint64())
}

func GetCurrentTime(this js.Value, args []js.Value) interface{} {
	return time.Now()
}
