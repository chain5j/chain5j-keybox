// Copyright 2014 The go-ethereum Authors
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
	"encoding/json"

	"github.com/pborman/uuid"
)

const (
	version = 3
)

type Key struct {
	Id   uuid.UUID // Version 4 "random" for unique id not derived from key data
	Path string    `json:"path"`
	Addr string
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey []byte
}

type plainKeyJSON struct {
	PrivateKey string `json:"privatekey"`
	Path       string `json:"path"`
	Id         string `json:"id"`
	Version    int    `json:"version"`
}

type encryptedKeyJSONV3 struct {
	Id      string     `json:"id"`
	Version int        `json:"version"`
	Address string     `json:"address,omitempty"`
	Crypto  CryptoJSON `json:"crypto"`
}

type encryptedKeyJSONV1 struct {
	Id      string     `json:"id"`
	Version string     `json:"version"`
	Crypto  CryptoJSON `json:"crypto"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

func (k *Key) MarshalJSON() (j []byte, err error) {
	jStruct := plainKeyJSON{
		hex.EncodeToString(k.PrivateKey),
		k.Path,
		k.Id.String(),
		version,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *Key) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSON)
	err = json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}

	u := new(uuid.UUID)
	*u = uuid.Parse(keyJSON.Id)
	k.Id = *u
	privkey, err := hex.DecodeString(keyJSON.PrivateKey)
	if err != nil {
		return err
	}

	k.PrivateKey = privkey
	k.Path = keyJSON.Path

	return nil
}
