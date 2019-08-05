package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// 随机生成UUID 的V4版本
func newUUIDV4String() (uuid string) {
	randBytes := newUuidV4()
	return _uuidToString(randBytes)
}

func newUuidV4() (randBytes [16]byte) {
	for k := range randBytes {
		randBytes[k] = byte(rand.Int() & 0xFF)
	}
	randBytes[6] &= 0xFF
	randBytes[6] |= 0x40
	randBytes[8] &= 0x3F
	randBytes[8] |= 0x80
	return randBytes
}
func _uuidToString(data [16]byte) (uuid string) {
	var msb, lsb int64
	msb = 0
	lsb = 0
	for i := 0; i < 8; i++ {
		msb = (msb << 8) | int64(data[i]&0xff)
	}
	for i := 8; i < 16; i++ {
		lsb = (lsb << 8) | int64(data[i]&0xff)
	}
	uuid = _digits(msb>>32, 8) + "-" +
		_digits(msb>>16, 4) + "-" +
		_digits(msb, 4) + "-" +
		_digits(lsb>>48, 4) + "-" +
		_digits(lsb, 12)
	return
}

func _digits(val int64, digits int) (hex string) {
	i := uint(digits * 4)
	var hi uint
	hi = 1 << i
	res := int64(hi) | val&int64(hi-1)
	format := "%0" + strconv.Itoa(digits) + "x"
	hex = fmt.Sprintf(format, res)
	hex = hex[1:]
	return

}
