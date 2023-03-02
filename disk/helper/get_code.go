package helper

import (
	"cloud-disk/disk/define"
	"math/rand"
	"time"
)

func GetEmailCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(10)])
	}
	return code
}
