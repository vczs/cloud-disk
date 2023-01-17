package helper

import uuid "github.com/satori/go.uuid"

func GetUid() string {
	return uuid.NewV4().String()
}
