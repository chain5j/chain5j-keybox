package p256

import (
	"crypto/elliptic"
	"fmt"

	"github.com/chain5j/chain5j-pkg/crypto/signature/prime256v1"
	"github.com/chain5j/keybox/algorithm"
)

type Algorithm struct {
}

// 从私钥中获取公钥
func (a *Algorithm) GetPubKeyFromPriKey(priKey []byte) ([]byte, error) {
	if len(priKey) == 0 {
		return nil, fmt.Errorf("Chain GetPubKeyFromPriKey parameter error")
	}
	_, publicKey := prime256v1.PrivKeyFromBytes(elliptic.P256(), priKey)
	return publicKey.SerializeUncompressed(), nil
}

// 签名交易体Hash
func (a *Algorithm) Sign(priKey []byte, hash []byte) (*algorithm.Signature, error) {
	privateKey, publicKey := prime256v1.PrivKeyFromBytes(elliptic.P256(), priKey)
	signBytes, err := prime256v1.SignCompact(privateKey, hash, false)
	if err != nil {
		return nil, err
	}
	return &algorithm.Signature{
		SignBytes: signBytes[0:64],
		V:         signBytes[64],
		Pubkey:    publicKey.SerializeUncompressed(),
	}, nil
}
