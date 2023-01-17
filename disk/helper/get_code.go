package helper

import (
	"cloud-disk/disk/internal/config"
	"math/rand"
	"time"
)

func GetEmailCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < config.CodeLength; i++ {
		code += string(s[rand.Intn(10)])
	}
	return code
}
