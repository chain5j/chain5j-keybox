// description: subchain-sdk-go
//
// @author: xwc1125
// @date: 2020/8/6 0006
package scrypt

import "errors"

var (
	ErrDecrypt = errors.New("could not decrypt key with given password")
)
