// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package eth

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"testing"
)

func TestChain_GetAddressFromPubKey(t *testing.T) {
	chain := NewChain("mainnet")
	privateKey, err := crypto.HexToECDSA(crypto.S256,"0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818")
	if err != nil {
		panic(err)
	}
	prvKeyBytes := crypto.FromECDSA(privateKey)
	pubBytes, err := chain.GetPubKeyFromPriKey(prvKeyBytes)
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
