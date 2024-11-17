// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package address

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

func TestBTCAddress(t *testing.T) {
	cryptoS256 := secp256k1.Secp251k1{}
	publicKey, _ := cryptoS256.GenerateKey(btcec.S256())
	pubKeyBytes, _ := cryptoS256.MarshalPublicKey(&publicKey.PublicKey)
	fmt.Println("pubKey", hexutil.Encode(pubKeyBytes))
	address := BTCAddress(BTCMainNet, pubKeyBytes)
	fmt.Println("address", address)
	ok := IsValidBTCAddress(address)
	fmt.Println(ok)
	publicKeyHashFromAddress := AddressToPubKeyRipemd160Hash(address)
	fmt.Println("pubKeyHash", hexutil.Encode(publicKeyHashFromAddress))
}
