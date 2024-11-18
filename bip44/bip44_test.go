package bip44

import (
	"encoding/json"
	"testing"

	"github.com/chain5j/keybox/bip32"
	"github.com/chain5j/keybox/bip39"
)

func Test_NewKeyFromMnemonic(t *testing.T) {
	seed, err := bip39.NewSeedWithErrorChecking("fragile disorder legal weapon depend sunny detail lens expect fresh dutch blur", "")
	if err != nil {
		t.Log(err)
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		t.Log(err)
	}

	key, err := json.Marshal(masterKey)
	if err != nil {
		t.Log(err)
	}

	t.Log(string(key))

	xkey, _ := NewKeyFromMasterKey(masterKey, 1, 0, 0, 0)
	t.Log(bip32.JsonString(xkey))

}

func Test_NewKeyFromMasterKey(t *testing.T) {

}
