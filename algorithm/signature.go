// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package algorithm

import "encoding/json"

type Signature struct {
	SignBytes []byte `json:"signBytes"` // 签名数据(R,S)
	V         byte   `json:"v"`         // 校验码
	Pubkey    []byte `json:"pubkey"`    // 公钥
}

func (s *Signature) VNone() []byte {
	return s.SignBytes
}

func (s *Signature) VLeft() []byte {
	return append([]byte{s.V}, s.SignBytes...)
}

func (s *Signature) VRight() []byte {
	return append(s.SignBytes, s.V)
}

func (s *Signature) SignWithPubkey() []byte {
	return append(s.SignBytes, s.Pubkey...)
}

func (s *Signature) Bytes() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}
