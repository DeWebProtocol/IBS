package hash

import (
	"crypto/sha1"
	"encoding/binary"
)

func Hash64(i uint64) uint64 {
	sha := sha1.New()
	sha.Write(uint64ToBytes(i))
	return bytesToInt64(sha.Sum(nil))
}
func uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}
func bytesToInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}
