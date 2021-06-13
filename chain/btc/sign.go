// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/24 0024
package btc

import (
	"errors"
	log "github.com/chain5j/log15"
)

func (c *Chain) SignRawTx(rawTx, privateKeyWif string) (signedRawTx string, err error) {
	msg := new(CustomHexMsg)
	err = msg.UnmarshalJSON(rawTx)
	if err != nil {
		return
	}
	msg.PrivKeys = &[]string{privateKeyWif}
	if msg.Flags == nil {
		var flagALL = "ALL"
		msg.Flags = &flagALL
	}
	signCmd := &SignRawTransactionCmd{
		RawTx:    msg.RawTx,
		Inputs:   msg.Inputs,
		PrivKeys: msg.PrivKeys,
		Flags:    msg.Flags,
	}
	conf, err := ParseNetworkToConf(string(c.networkType))
	result, err := SignRawTransaction(signCmd, conf)
	if err != nil {
		return
	}
	if result.Errors != nil && len(result.Errors) > 0 {
		log.Error("BTC SignRawTransaction err", "err", result.Errors)
		errs := ""
		for _, e := range result.Errors {
			errs = errs + e.Error
		}
		err = errors.New(errs)
		return
	}
	signedRawTx = result.Hex
	return
}
