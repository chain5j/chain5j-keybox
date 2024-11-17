// description: keybox
//
// @author: xwc1125
// @date: 2020/8/6 0006
package keybox

import "github.com/chain5j/keybox/algorithm"

type ChainAPI interface {
	algorithm.AlgorithmAPI
	ChainInfo() *ChainInfo                                                 // 链内容
	ExportPrivateKey(priKey []byte, isCompressPubKey bool) (string, error) // 导出私钥
	GetAddressFromPubKey(pubKey []byte) (string, error)                    // 通过公钥获取地址
	SignToStr(priKey []byte, hash []byte) (string, error)                  // 签名直接返回签名的string
}
