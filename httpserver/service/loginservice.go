package service

import (
	"crypto/sha1"
	"encoding/hex"
	"redissyncer-portal/httpserver/model"
	"strconv"
	"time"
)

func LoginService(body model.Login) (string, error) {

	myhash := sha1.New()
	myhash.Write([]byte(body.User + body.Password + strconv.Itoa(time.Now().Nanosecond())))
	bs := myhash.Sum(nil)
	laststr := hex.EncodeToString(bs)
	return laststr, nil
}
