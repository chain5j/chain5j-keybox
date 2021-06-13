// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package address

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"testing"
)

func TestBTCAddress(t *testing.T) {
	publicKey, _ := crypto.GenerateKey(crypto.S256)
	pubKeyBytes := crypto.MarshalPubkey(&publicKey.PublicKey)
	fmt.Println("pubKey", hexutil.Encode(pubKeyBytes))
	address := BTCAddress(BTCMainNet, pubKeyBytes)
	fmt.Println("address", address)
	ok := IsValidBTCAddress(address)
	fmt.Println(ok)
	publicKeyHashFromAddress := AddressToPubKeyRipemd160Hash(address)
	fmt.Println("pubKeyHash", hexutil.Encode(publicKeyHashFromAddress))
}
