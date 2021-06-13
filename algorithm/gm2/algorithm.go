package gm2

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto/gmsm"
	"github.com/chain5j/keybox/algorithm"
	"github.com/tjfoc/gmsm/sm2"
)

type Algorithm struct {
}

// 从私钥中获取公钥
func (a *Algorithm) GetPubKeyFromPriKey(priKey []byte) ([]byte, error) {
	if len(priKey) == 0 {
		return nil, fmt.Errorf("Chain GetPubKeyFromPriKey parameter error")
	}
	_, publicKey := gmsm.PrivKeyFromBytes(priKey)
	return sm2.Compress(publicKey), nil
}

// 签名交易体Hash
func (a *Algorithm) Sign(priKey []byte, hash []byte) (*algorithm.Signature, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	privateKey, publicKey := gmsm.PrivKeyFromBytes(priKey)
	r, b, err := sm2.Sign(privateKey, hash)
	if err != nil {
		return nil, err
	}
	bytes, err := sm2.SignDigitToSignData(r, b)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := sm2.Compress(publicKey)
	if err != nil {
		return nil, err
	}
	return &algorithm.Signature{
		SignBytes: bytes,
		V:         0, // 国密不支持通过签名内容恢复公钥
		Pubkey:    pubKeyBytes,
	}, nil
}
