package s256

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/keybox/algorithm"
)

type Algorithm struct {
}

// 从私钥中获取公钥
func (a *Algorithm) GetPubKeyFromPriKey(priKey []byte) ([]byte, error) {
	if len(priKey) == 0 {
		return nil, fmt.Errorf("Chain GetPubKeyFromPriKey parameter error")
	}
	privateKey, err := crypto.ToECDSA(crypto.S256, priKey)
	if err != nil {
		return nil, err
	}
	return crypto.MarshalPubkey(&privateKey.PublicKey), nil
}

// 签名交易体Hash
func (a *Algorithm) Sign(priKey []byte, hash []byte) (*algorithm.Signature, error) {
	privateKey, err := crypto.ToECDSA(crypto.S256, priKey)
	if err != nil {
		return nil, err
	}
	signResult, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	return &algorithm.Signature{
		SignBytes: signResult.Signature[:64],
		V:         signResult.Signature[64],
		Pubkey:    crypto.MarshalPubkey(&privateKey.PublicKey),
	}, nil
}
