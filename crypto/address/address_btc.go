// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package address

import (
	"bytes"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	mainnetVersion     = byte(0x00) // 定义版本号，一个字节[mainnet]
	testnetVersion     = byte(0x6f) // [regtest]
	testnet3Version    = byte(0x6f) // [testnet3]
	simnetVersion      = byte(0x3f) // [simnet]
	addressChecksumLen = 4          // 定义checksum长度为四个字节
)

type BTCNetType string

const (
	BTCMainNet  BTCNetType = "mainnet"
	BTCTestNet  BTCNetType = "regtest"
	BTCTestNet3 BTCNetType = "testnet3"
	BTCSimNet   BTCNetType = "simnet"
)

// pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
func BTCAddress(netType BTCNetType, pubKey []byte) string {
	// 调用Ripemd160Hash返回160位的Pub Key hash
	ripemd160Hash := Ripemd160Hash(pubKey)

	// 将version+Pub Key hash
	var version byte
	switch netType {
	case BTCMainNet:
		version = mainnetVersion
	case BTCTestNet:
		version = testnetVersion
	case BTCTestNet3:
		version = testnet3Version
	case BTCSimNet:
		version = simnetVersion
	default:
		version = mainnetVersion
	}
	versionRipemd160hash := append([]byte{version}, ripemd160Hash...)

	// 调用CheckSum方法返回前四个字节的checksum
	checkSumBytes := CheckSum(versionRipemd160hash)

	// 将version+Pub Key hash+ checksum生成25个字节
	bytes := append(versionRipemd160hash, checkSumBytes...)

	// 将这25个字节进行base58编码并返回
	return base58.Encode(bytes)
}

// sha256(sha256(versionPublickeyHash))  取最后4个字节的值
func CheckSum(payload []byte) []byte {
	// 这里传入的payload其实是version+Pub Key hash，对其进行两次256运算
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen] // 返回前四个字节，为CheckSum值
}

// 对公钥进行sha256散列和ripemd160散列,获得publickeyHash
func Ripemd160Hash(publicKey []byte) []byte {
	// 将传入的公钥进行256运算，返回256位hash值
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	// 将上面的256位hash值进行160运算，返回160位的hash值
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil) // 返回Pub Key hash
}

// 通过地址获得公钥的ripemd160Hash
func AddressToPubKeyRipemd160Hash(address string) []byte {
	fullHash := base58.Decode(address)
	publicKeyHash := fullHash[1 : len(fullHash)-addressChecksumLen]
	return publicKeyHash
}

// 判断地址是否有效
func IsValidBTCAddress(address string) bool {
	// 将地址进行base58反编码，生成的其实是version+Pub Key hash+ checksum这25个字节
	versionPublicCheckSumBytes := base58.Decode(address)

	// [25-4:],就是21个字节往后的数（22,23,24,25一共4个字节）
	checkSumBytes := versionPublicCheckSumBytes[len(versionPublicCheckSumBytes)-addressChecksumLen:]
	// [:25-4],就是前21个字节（1～21,一共21个字节）
	versionRipemd160 := versionPublicCheckSumBytes[:len(versionPublicCheckSumBytes)-addressChecksumLen]
	// 取version+public+checksum的字节数组的前21个字节进行两次256哈希运算，取结果值的前4个字节
	checkBytes := CheckSum(versionRipemd160)
	// 将checksum比较，如果一致则说明地址有效，返回true
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}

	return false
}
