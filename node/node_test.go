package node

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

func TestSha1(t *testing.T) {
	sha := sha1.New()
	sha.Write([]byte{'1'})
	hash := sha.Sum(nil)

	fmt.Println(hex.EncodeToString(hash))
}
func TestMap(t *testing.T) {
	m := make(map[int]int)
	m[0] = 1
	b, ok := m[0]
	fmt.Println(b, ok)
}
func TestCorruptFunctionLinear(t *testing.T) {
	acc := 0
	for i := 0; i < 10000; i++ {
		r := rand.Intn(10000)
		if r <= i {
			acc++
		}
	}
	fmt.Println(acc)
}
