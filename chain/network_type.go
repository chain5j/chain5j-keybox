// description: keybox
//
// @author: xwc1125
// @date: 2020/8/20 0020
package chain

type NetworkType string

const (
	MainNet = "mainnet"
	TestNet = "testnet"
	DevNet  = "devnet"
)

func ParseToType(network string) NetworkType {
	switch network {
	case MainNet:
		return MainNet
	case TestNet:
		return TestNet
	case DevNet:
		return DevNet
	default:
		return MainNet
	}
}
