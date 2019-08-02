package main

import (
	"encoding/base64"
)

func decodeBase64(str string) (res []byte) {
	if str == "" {
		return
	}
	res, e := base64.StdEncoding.DecodeString(str)
	if e != nil {
		return
	}
	return
}
