// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package btc

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

func TestChain_GetAddressFromPubKey(t *testing.T) {
	pubBytes, err := hexutil.Decode("0x042ed4e6b73adb88f12cfc8ec7fd2c8423c9577857e6a6e37226721967d995368def472a77c125c780767a5eb28e443b7591bfb6fd134d3b47624f17194d0a5993")
	if err != nil {
		panic(err)
	}
	chain := NewChain("mainnet")
	addr, err := chain.GetAddressFromPubKey(pubBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("addr", addr)
}

func TestPrvKey(t *testing.T) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
	prvKeyBytes := privKey.Serialize()
	fmt.Println("privKey", hexutil.Encode(prvKeyBytes))
	pubKeyBytes := privKey.PubKey().SerializeUncompressed()
	fmt.Println("pubKey", hexutil.Encode(pubKeyBytes))

	fmt.Println("========================")

	chain := NewChain("mainnet")
	addr, err := chain.GetAddressFromPubKey(pubKeyBytes)
	fmt.Println("keybox addr", addr)
	wifPrvkey, err := chain.ExportPrivateKey(prvKeyBytes, false)
	if err != nil {
		panic(err)
	}
	fmt.Println("keybox privKeyWif", wifPrvkey)

	fmt.Println("========================")
	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("btc addr", pubKeyAddress.EncodeAddress())
	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	if err != nil {
		panic(err)
	}
	fmt.Println("privKeyWif", privKeyWif.String())
}

func GenerateBTC() (string, string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", err
	}

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", "", err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}

	return privKeyWif.String(), pubKeyAddress.EncodeAddress(), nil
}

func GenerateBTCTest() (string, string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", err
	}

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, false)
	if err != nil {
		return "", "", err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.TestNet3Params)
	if err != nil {
		return "", "", err
	}

	return privKeyWif.String(), pubKeyAddress.EncodeAddress(), nil
}

func TestBtcGeneTest(t *testing.T) {
	wifKey, address, _ := GenerateBTCTest() // 测试地址
	// wifKey, address, _ := GenerateBTC() // 正式地址
	fmt.Println("address", address)
	fmt.Println("wifKey", wifKey)
}

func TestBtcGeneProd(t *testing.T) {
	wifKey, address, _ := GenerateBTC() // 正式地址
	fmt.Println("address", address)
	fmt.Println("wifKey", wifKey)
}

func TestWifToAddr(t *testing.T) {
	// cW2gNjzkXbcgHJrus1A99cW8J3STUTfSkvvjaU3r42ayAsntZiwJ
	// [RegressionNetParams] mvmSUX991W3GrhYzjQX84qWduQDE8EBnfW
	wifStr := "cW2gNjzkXbcgHJrus1A99cW8J3STUTfSkvvjaU3r42ayAsntZiwJ"
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		panic(err)
	}
	pubKeySerial := wif.SerializePubKey()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	println("addr", pubKeyAddress.EncodeAddress())
}
