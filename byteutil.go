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

func encodeBase64(data []byte) (res string) {
	if data == nil {
		return
	}
	res = base64.StdEncoding.EncodeToString(data)
	return
}
