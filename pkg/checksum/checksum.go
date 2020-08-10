package checksum

import (
	"crypto/sha1"
	"fmt"
)

//ComputeSha1 takes a slice of bytes and returns the SHA1 hash of them
func ComputeSha1(str []byte) string {
	//sha1.Sum() returns an array of 20 bytes
	//fmt.Printf("% x", byteArr) will show each byte -- represented as
	//hex digits -- with a space between, as in
	//'30 85 29 60 69 76 2c 17 b3 6f 0c b5 db 81 10 c6 54 b4 d6 69'
	//(for our purposes, no inter-byte spacing necessary)
	byteArr := sha1.Sum(str)
	return fmt.Sprintf("%x", byteArr)
}
