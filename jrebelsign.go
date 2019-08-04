package main

import (
	"strconv"
)

func toLeaseCreateJson(clientRandomness string, serverRandomness string, guid string, offline bool, validFrom string, validUntil string) (res string) {
	var s2 string
	if offline {
		s2 = clientRandomness + ";" + serverRandomness + ";" + guid + ";" + strconv.FormatBool(offline) + ";" + validFrom + ";" + validUntil
	} else {
		s2 = clientRandomness + ";" + serverRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
	}
	signature, err := signWithSha1([]byte(s2))
	if err != nil {
		return
	}
	res = encodeBase64(signature)
	return
}
