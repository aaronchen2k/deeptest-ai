package _str

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func Uuid() string {
	uid, _ := uuid.NewV4()
	return strings.Replace(uid.String(), "-", "", -1)
}

func UuidWithSep() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}

func Ulid() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	rand, _ := ulid.New(ms, entropy)

	ret := strings.ToLower(rand.String())
	ret = strings.Replace(ret, "-", "", -1)

	return ret
}
