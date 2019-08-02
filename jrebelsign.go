package main

import (
	"strconv"
)

//服务端随机数,如果要自己生成，务必将其写到json的serverRandomness中
const serverRandomness string = "H2ulzLlh7E0="

func toLeaseCreateJson(clientRandomness string, guid string, offline bool, validFrom string, validUntil string) string {
	var s2 string
	if offline {
		s2 = clientRandomness + ";" + serverRandomness + ";" + guid + ";" + strconv.FormatBool(offline) + ";" + validFrom + ";" + validUntil
	} else {
		s2 = clientRandomness + ";" + serverRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
	}
	return s2
}
