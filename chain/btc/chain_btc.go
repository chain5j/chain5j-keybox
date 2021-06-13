package btc

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/keybox"
	"github.com/chain5j/keybox/algorithm/s256"
	"github.com/chain5j/keybox/chain"
	"github.com/chain5j/keybox/crypto/address"
)

type Chain struct {
	s256.Algorithm
	chainInfo   *keybox.ChainInfo
	networkType chain.NetworkType
	netId       byte
}

func NewChain(networkType chain.NetworkType) *Chain {
	netId := byte(0x80)
	switch networkType {
	case chain.MainNet:
		netId = byte(0x80)
	case chain.TestNet:
		netId = byte(0xef)
	case chain.DevNet:
		netId = byte(0xef)
	}
	return &Chain{
		chainInfo:   keybox.ChainInfoBTC,
		networkType: networkType,
		netId:       netId,
	}
}

// 获取链信息
func (a *Chain) ChainInfo() *keybox.ChainInfo {
	return a.chainInfo
}

// 比特币中的netId就是PrivateKeyID[wifPrvkey]
func (a *Chain) ExportPrivateKey(priKey []byte, isCompressPubKey bool) (string, error) {
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), priKey)
	parseNetworkToConf, err := ParseNetworkToConf(string(a.networkType))
	if err != nil {
		return "", err
	}
	privKeyWif, err := btcutil.NewWIF(privateKey, parseNetworkToConf, isCompressPubKey)
	if err != nil {
		return "", err
	}
	return privKeyWif.String(), nil
}

// 从公钥获取地址
func (a *Chain) GetAddressFromPubKey(pubKey []byte) (string, error) {
	if pubKey == nil || len(pubKey) == 0 {
		return "", fmt.Errorf("pubKey is empty")
	}
	var netType = address.BTCMainNet
	switch a.networkType {
	case chain.MainNet:
		netType = address.BTCMainNet
	case chain.TestNet:
		netType = address.BTCTestNet
	case chain.DevNet:
		netType = address.BTCTestNet3
	}
	return address.BTCAddress(netType, pubKey), nil
}

// 签名直接返回签名的string
func (a *Chain) SignToStr(priKey []byte, rawTxBytes []byte) (string, error) {
	wifPrvkey, err := a.ExportPrivateKey(priKey, false)
	fmt.Println("wifPrvkey", wifPrvkey)
	if err != nil {
		return "", err
	}
	signedRawTx, err := a.SignRawTx(hex.EncodeToString(rawTxBytes), wifPrvkey)
	if err != nil {
		return "", err
	}
	return signedRawTx, nil
}

const compressMagic byte = 0x01

// 比特币中的netId就是PrivateKeyID[wifPrvkey]
func (a *Chain) ExportPrivateKey2(priKey []byte, isCompressPubKey bool) (string, error) {
	encodeLen := 1 + btcec.PrivKeyBytesLen + 4
	if isCompressPubKey {
		encodeLen++
	}
	privateKey, err := crypto.ToECDSA(crypto.S256, priKey)
	if err != nil {
		return "", err
	}

	p := make([]byte, 0, encodeLen)
	p = append(p, a.netId)
	p = paddedAppend(btcec.PrivKeyBytesLen, p, privateKey.D.Bytes())
	if isCompressPubKey {
		p = append(p, compressMagic)
	}
	cksum := chainhash.DoubleHashB(p)[:4]
	p = append(p, cksum...)
	return base58.Encode(p), nil
}

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}
