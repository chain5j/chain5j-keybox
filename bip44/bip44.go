package bip44

// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

import (
	"github.com/chain5j/keybox/bip32"
	"github.com/chain5j/keybox/bip39"
)

const CoinTypeBTC uint32 = 0x80000000
const Purpose uint32 = 0x8000002C

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
// const (
//	TypeBTC  uint32 = 0x80000000
//	TypeETH  uint32 = 0x8000003c
//	TypeOMNI uint32 = 0x800000c8
// )

func NewKeyFromMnemonic(mnemonic string, coin, account, chain, address uint32) (*bip32.Key, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	return NewKeyFromMasterKey(masterKey, coin, account, chain, address)
}

// m / purpose' / coin_type' / account' / change / address_index
// purpose 根据BIP43建议将常量设置为44'（或0x8000002C）。它指示根据此规范使用了此节点的子树
// coin_type 特指币种并且允许多元货币 HD 钱包中的货币在第二个层级下有自己的亚树状结构
// account  将密钥空间划分为独立的用户身份
// change 0用于外部接收地址 1用于找零地址
// address_index  地址索引

func NewKeyFromMasterKey(masterKey *bip32.Key, coin_type, account, change, address_index uint32) (*bip32.Key, error) {
	child, err := masterKey.NewChildKey(Purpose)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(coin_type)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(account)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(address_index)
	if err != nil {
		return nil, err
	}

	return child, nil
}

func NewKeyFromMasterKeyWithOrg(masterKey *bip32.Key, purpose, coinType, org, account, change, addressIndex uint32) (*bip32.Key, error) {
	child, err := masterKey.NewChildKey(purpose)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(coinType)
	if err != nil {
		return nil, err
	}

	if purpose != Purpose {
		child, err = child.NewChildKey(org)
		if err != nil {
			return nil, err
		}
	}

	child, err = child.NewChildKey(account)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(addressIndex)
	if err != nil {
		return nil, err
	}

	return child, nil
}
