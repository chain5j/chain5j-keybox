// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package keybox

import "github.com/chain5j/keybox/bip44"

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	TypeBTC  uint32 = 0x80000000
	TypeETH  uint32 = 0x8000003c
	TypeOMNI uint32 = 0x800000c8
)

const (
	Purpose45 uint32 = 0x8000002d
)

type ChainInfo struct {
	ChainName     string // 链名称（eth，btc）
	ChainType     uint32 // 链分配的类型值
	AlgorithmName string // 链的算法名称（s256,p256,gm2）
	Algorithm     uint32 // 链算法类型值
}

var (
	ChainInfoBTC = &ChainInfo{
		ChainName:     "BTC",
		ChainType:     bip44.CoinTypeBTC,
		AlgorithmName: "S256",
		Algorithm:     0x80000200,
	}
	ChainInfoETH = &ChainInfo{
		ChainName:     "ETH",
		ChainType:     TypeETH,
		AlgorithmName: "S256",
		Algorithm:     0x80000200,
	}
)
