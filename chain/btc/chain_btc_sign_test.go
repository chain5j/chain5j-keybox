// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/20 0020
package btc

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/keybox/crypto/address"
	"testing"
)

func TestTransaction(t *testing.T) {
	/**
	{
	    "txid": "d1681e7b2eb6841457a80d53973237b29aa09b0f4bff8d121d7c6ae6446a99a9",
	    "vout": 1,
	    "address": "2Mx9f29Y4tmTztQtSKqb5Wgdjf7vqUacZ3X",
	    "label": "",
	    "scriptPubKey": "a91435cb0767a187e17a02eb7799045337229be7868087",
	    "amount": 23.30000000,
	    "confirmations": 0,
	    "spendable": false,
	    "solvable": false,
	    "safe": true
	  }
	*/
	// 1)添加unspent内容
	input := new(BTCUnspent)
	input.Add(
		"9c1193275242a1dfeb9cf1214af3252fab8281e3e54b2cb26de76db9e6d7ebff",
		0,
		12.5,
		"2103bd0d9c8e74b846f96e8a0beeeaa3a3e2785ff67cdf398bef7bf8824d7a07d4baac",
		"",
	)

	// 2)对output进行添加
	output := new(BTCOutput)
	toAddr, _ := NewBTCAddressFromString("myxu5JjH9zU5L2GEhaqiCUUjKm71SZ1hzp", "testnet")
	toAmount, _ := NewBTCAmount(10.0)
	output.Add(toAddr, toAmount)

	// 3)找零
	changeAddr, _ := NewBTCAddressFromString("mvmSUX991W3GrhYzjQX84qWduQDE8EBnfW", "testnet")

	// 4)拼接交易
	tt, err := NewBTCTransaction(input, output, changeAddr, 2, "testnet")
	if err != nil {
		t.Fatal(err)
	}

	ii, _ := tt.Encode()
	t.Log(ii)

	hh, err := tt.EncodeToSignCmd()
	if err != nil {
		t.Fatal(err)
	}
	// TODO 此处完成createrawtransaction
	t.Log("rawTx", hh)
	t.Log(tt.GetFee())

	chain := NewChain("testnet")
	// addr: msKe45XX3Sf6bnYM6UXpxbzpf4STqSFDkU    priKeyWif: 92gNnLAujLaLS1o7xtGnf8Xx4T25ZjaVccH2gXwECoPNHRgwuv1
	// addr: myxu5JjH9zU5L2GEhaqiCUUjKm71SZ1hzp    priKeyWif: cVBp35B945nC4AEHgAdLJQaGewuFJH4PXAgETBxRmmjavJZtQCAB
	// addr: mvmSUX991W3GrhYzjQX84qWduQDE8EBnfW    priKeyWif: cW2gNjzkXbcgHJrus1A99cW8J3STUTfSkvvjaU3r42ayAsntZiwJ
	signedRawTx, err := chain.SignRawTx(hh, "cW2gNjzkXbcgHJrus1A99cW8J3STUTfSkvvjaU3r42ayAsntZiwJ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signedRawTx", signedRawTx)
}

func TestChain_SignToStr(t *testing.T) {
	privateKeyWif := "cVBp35B945nC4AEHgAdLJQaGewuFJH4PXAgETBxRmmjavJZtQCAB"
	fmt.Println("privateKeyWif", privateKeyWif)
	wif, err := btcutil.DecodeWIF(privateKeyWif)
	if err != nil {
		panic(err)
	}
	btcAddress := address.BTCAddress(address.BTCTestNet, wif.SerializePubKey())
	fmt.Println("btcAddress", btcAddress)

	prvKeyBytes := wif.PrivKey.Serialize()
	chain := NewChain("testnet")
	privateKeyWif1, err := chain.ExportPrivateKey(prvKeyBytes, wif.CompressPubKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("privateKeyWif1", privateKeyWif1)
	privateKeyWif2, err := chain.ExportPrivateKey2(prvKeyBytes, wif.CompressPubKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("privateKeyWif2", privateKeyWif2)

	bytes, err := hexutil.Decode("7b225261775478223a2230313030303030303031643735613661333230376562613632616261633833326661313534653061323163303232353330346665303037613830356463346132643866313739373564363031303030303030303066666666666666663032303030303030303030303030303030303139373661393134383137646235303066656465643065323535363864356635333537633962636233316462313539343838616363303261636331643030303030303030313937366139313438313764623530306665646564306532353536386435663533353763396263623331646231353934383861633030303030303030222c22496e70757473223a5b7b2274786964223a2264363735373966316438613263343564383037613030666530343533323263303231306134653135666133326338626132616136656230373332366135616437222c22766f7574223a312c227363726970745075624b6579223a223736613931343831376462353030666564656430653235353638643566353335376339626362333164623135393438386163222c2272656465656d536372697074223a22227d5d2c22507269764b657973223a6e756c6c2c22466c616773223a6e756c6c7d")
	if err != nil {
		panic(err)
	}
	// 需要将isCompressPubKey改成false
	signToStr, err := chain.SignToStr(prvKeyBytes, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("signToStr", signToStr)
}
