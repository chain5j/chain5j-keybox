// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package scrypt

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

// Tests that a json key file can be decrypted and encrypted in multiple rounds.
func TestKeyEncryptDecrypt(t *testing.T) {
	keyjson, err := ioutil.ReadFile("testdata/very-light-scrypt.json")
	if err != nil {
		t.Fatal(err)
	}
	password := ""
	//address,_ := hex.DecodeString("0x45dea0fb0bba44f4fcf290bba71fd57d7117cbb8")

	// Do a few rounds of decryption and encryption
	for i := 0; i < 3; i++ {
		// Try a bad password first
		if _, err := DecryptKey(keyjson, password+"bad"); err == nil {
			t.Errorf("test %d: json key decrypted with bad password", i)
		}
		// Decrypt with the correct password
		key, err := DecryptKey(keyjson, password)
		if err != nil {
			t.Fatalf("test %d: json key failed to decrypt: %v", i, err)
		}
		//if key.Address != address {
		//	t.Errorf("test %d: key address mismatch: have %x, want %x", i, key.Address, address)
		//}
		// Recrypt with a new password and start over
		password += "new data appended"
		if keyjson, err = EncryptKey(key, password, veryLightScryptN, veryLightScryptP); err != nil {
			t.Errorf("test %d: failed to recrypt key %v", i, err)
		}
	}
}

func TestEncryptKey(t *testing.T) {
	prvKey, _ := hex.DecodeString("0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818")
	k := &Key{
		//Address:    addr,
		PrivateKey: prvKey,
	}
	keyJson, err := EncryptKey(k, "123456", veryLightScryptN, veryLightScryptP)
	if err != nil {
		t.Errorf("test: failed to recrypt key %v", err)
	}

	key, err := DecryptKey(keyJson, "123456")
	if err != nil {
		t.Errorf("test: failed to recrypt key %v", err)
	}
	prvKey1 := hex.EncodeToString(key.PrivateKey)
	fmt.Println("prvKey1", prvKey1)
}
