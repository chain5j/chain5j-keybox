// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/21 0021
package btc

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type BTCAddress struct {
	address btcutil.Address
}

// 将地址字符串转为BTC地址
func NewBTCAddressFromString(addr string, network string) (address *BTCAddress, err error) {
	netParams, err := ParseNetworkToConf(network)
	if err != nil {
		return nil, err
	}

	address = new(BTCAddress)
	decAddr, err := btcutil.DecodeAddress(addr, netParams)
	if err != nil {
		return
	}
	address.address = decAddr
	return
}

func ParseNetworkToConf(network string) (*chaincfg.Params, error) {
	switch network {
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	case "testnet3":
		return &chaincfg.TestNet3Params, nil
	case "testnet":
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, fmt.Errorf("network is error")
	}
}

// BTCAmount 对float64 的封装
type BTCAmount struct {
	amount btcutil.Amount
}

// NewBTCAmount 数量in BTC (not in satoshi)
func NewBTCAmount(amount float64) (amt *BTCAmount, err error) {
	amt = new(BTCAmount)
	tempAmt, err := btcutil.NewAmount(amount)
	if err != nil {
		return
	}
	amt.amount = tempAmt
	return
}
