// description: keybox
//
// @author: xwc1125
// @date: 2020/8/6 0006
package scrypt

import "golang.org/x/crypto/sha3"

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
