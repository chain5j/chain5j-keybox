// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package algorithm

type AlgorithmAPI interface {
	GetPubKeyFromPriKey(priKey []byte) ([]byte, error)   // 通过私钥获取公钥
	Sign(priKey []byte, hash []byte) (*Signature, error) // 使用私钥对交易体Hash进行签名
}
