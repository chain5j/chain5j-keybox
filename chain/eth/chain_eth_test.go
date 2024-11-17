// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package eth

import (
	"fmt"
	"testing"

	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

func TestChain_GetAddressFromPubKey(t *testing.T) {
	chain := NewChain("mainnet")
	// cryptoS256 := secp256k1.Secp251k1{}
	priKey, _ := hexutil.Decode("0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818")
	// privateKey1, _ := btcec.PrivKeyFromBytes(priKey)
	// privateKey := cryptoS256.ToECDSA(privateKey1)
	// prvKeyBytes := cryptoS256.FromECDSA(privateKey)

	pubBytes, err := chain.GetPubKeyFromPriKey(priKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("pubKey", hexutil.Encode(pubBytes))

	addr, err := chain.GetAddressFromPubKey(pubBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("addr", addr)
}
