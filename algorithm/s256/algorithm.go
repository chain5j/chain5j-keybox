package s256

import (
	"fmt"

	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"github.com/chain5j/keybox/algorithm"
)

type Algorithm struct {
}

// 从私钥中获取公钥
func (a *Algorithm) GetPubKeyFromPriKey(priKey []byte) ([]byte, error) {
	if len(priKey) == 0 {
		return nil, fmt.Errorf("Chain GetPubKeyFromPriKey parameter error")
	}
	cryptoS256 := secp256k1.Secp251k1{}
	privateKey := cryptoS256.ToECDSA(priKey)
	if privateKey == nil {
		return nil, fmt.Errorf("private key is empty")
	}
	return cryptoS256.MarshalPublicKey(&privateKey.PublicKey)
}

// 签名交易体Hash
func (a *Algorithm) Sign(priKey []byte, hash []byte) (*algorithm.Signature, error) {
	cryptoS256 := secp256k1.Secp251k1{}
	privateKey := cryptoS256.ToECDSA(priKey)
	if privateKey != nil {
		return nil, fmt.Errorf("private key is empty")
	}
	signResult, err := cryptoS256.Sign(privateKey, hash)
	if err != nil {
		return nil, err
	}
	marshalPublicKey, err := cryptoS256.MarshalPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return &algorithm.Signature{
		SignBytes: signResult[:64],
		V:         signResult[64],
		Pubkey:    marshalPublicKey,
	}, nil
}
