package eth

import (
	"encoding/hex"
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto/keccak"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/keybox"
	"github.com/chain5j/keybox/algorithm/s256"
	"github.com/chain5j/keybox/chain"
)

type Chain struct {
	s256.Algorithm
	chainInfo   *keybox.ChainInfo
	networkType chain.NetworkType
}

func NewChain(networkType chain.NetworkType) *Chain {
	return &Chain{
		chainInfo:   keybox.ChainInfoETH,
		networkType: networkType,
	}
}

// 获取链信息
func (a *Chain) ChainInfo() *keybox.ChainInfo {
	return a.chainInfo
}

// 导出私钥
func (a *Chain) ExportPrivateKey(priKey []byte, isCompressPubKey bool) (string, error) {
	return hexutil.Encode(priKey), nil
}

// 从公钥获取地址
func (a *Chain) GetAddressFromPubKey(pubKey []byte) (string, error) {
	if pubKey == nil || len(pubKey) == 0 {
		return "", fmt.Errorf("pubKey is empty")
	}
	bytes := keccak.Keccak256(pubKey[1:])[12:]
	return types.BytesToAddress(bytes).Hex(), nil
}

// 签名直接返回签名的string
func (a *Chain) SignToStr(priKey []byte, hash []byte) (string, error) {
	signature, err := a.Sign(priKey, hash)
	if err != nil {
		return "", err
	}
	signBytes := signature.VRight()
	return hex.EncodeToString(signBytes), nil
}
