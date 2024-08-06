/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: address.go
 * Date: 2018-08-31
 *
 */

package ligo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/tyler-smith/go-bip39"
)

const (
	VersionNormalAddr = 0x0
)

/**
 * LIGO address struct
 */
type LigoAddr struct {
	Addr string //ligo main chain address string
}

// using BIP44 to manage ligo address
// m / purpose' / coin_type' / account' / change / address_index
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

func MnemonicToSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}

func getMasterkey(seed []byte, mainnet bool) (*hdkeychain.ExtendedKey, error) {

	// main net
	if mainnet {

		return hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	} else { //test net

		return hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)

	}

}
func GetNewPrivate() (privWif string, pubWif string, addr string, err error) {
	priv, err := btcec.NewPrivateKey()
	if err != nil {
		return
	}
	tmp_Wif, err := btcutil.NewWIF(priv, &chaincfg.MainNetParams, true)
	if err != nil {
		return
	}
	privWif = tmp_Wif.String()

	pubWif, _ = DerivePubkey(privWif)
	addr, err = getLigoAddressByWif(privWif, VersionNormalAddr)
	if err != nil {
		return
	}
	return
}
func getLigoAddressByWif(wif string, version uint32) (addr string, err error) {
	addr = ""
	wif_key, err := getPrivKey(wif)
	if err != nil {
		return
	}
	addr = GetAddressByPubkey(wif_key.PubKey().SerializeCompressed(), "main", version)

	return
}

func getAccountExtentkey(masterKey *hdkeychain.ExtendedKey, account uint32, addrIndex uint32) (*hdkeychain.ExtendedKey, string, error) {
	path := fmt.Sprintf("")
	// m / purpose’ / coin’ / account’ / change / address_index
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	// purpose & coin_type & change
	purpose := uint32(0x8000002C)
	coinType := uint32(0x80000000)
	accountType := uint32(0x80000000 + account)
	change := uint32(0)

	// m / 44'
	//masterKey.Derive()
	purposeKey, err := masterKey.Derive(purpose)
	path = path + "m/44'"
	if err != nil {
		return nil, path, fmt.Errorf("create purpose key failed: %v", err)
	}

	// m / 44' / 0'
	coinTypeKey, err := purposeKey.Derive(coinType)
	path = path + "/0'"
	if err != nil {
		return nil, path, fmt.Errorf("create coin type key failed: %v", err)
	}

	// m / 44' / 0' / 0'
	accountTypeKey, err := coinTypeKey.Derive(accountType)
	path = path + "/0'"
	if err != nil {
		return nil, path, fmt.Errorf("create account type key failed: %v", err)
	}

	// m / 44' / 0' / 0' / change
	changeKey, err := accountTypeKey.Derive(change)
	path = fmt.Sprintf("%s/%d", path, change)
	if err != nil {
		return nil, path, fmt.Errorf("create change key failed: %v", err)
	}

	// m / 44' / 0' / 0' / change / index
	addrIndexKey, err := changeKey.Derive(addrIndex)
	path = fmt.Sprintf("%s/%d", path, addrIndex)
	if err != nil {
		return nil, path, fmt.Errorf("create addr index key failed: %v", err)
	}

	return addrIndexKey, path, err
}

func getAccountAddr(addressKey *hdkeychain.ExtendedKey, nettype string, version uint32) (string, error) {

	ecPubkey, err := addressKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("get ecPubkey failed: %v", err)
	}

	// sha512 & ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	return GetAddressByPubkey(pubkeyByte, nettype, version), nil
}

func GetAddressByPubkey(pubkeyByte []byte, nettype string, version uint32) string {
	sha256Byte := sha256.Sum256(pubkeyByte)
	myRipemd := ripemd160.New()
	myRipemd.Write(sha256Byte[:])
	addrByte := myRipemd.Sum(nil)

	calcByte := []byte{byte(version)}
	calcByte = append(calcByte, addrByte...)
	addrByteChecksum := sha256.Sum256(calcByte)
	addrByteChecksum = sha256.Sum256(addrByteChecksum[:])
	calcByte = append(calcByte, addrByteChecksum[0:4]...)

	return base58.Encode(calcByte)
}

func ValidateAddress(address, net string) bool {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	buf := base58.Decode(address)
	if len(buf) != 25 {
		return false
	}

	tocheck := buf[0:21]
	check := buf[21:25]

	addrByteChecksum := sha256.Sum256(tocheck)
	addrByteChecksum = sha256.Sum256(addrByteChecksum[:])

	return bytes.Compare(check, addrByteChecksum[0:4]) == 0
}

func getWifkey(addressKey *hdkeychain.ExtendedKey) (string, error) {

	ecPrivkey, err := addressKey.ECPrivKey()
	if err != nil {
		return "", fmt.Errorf("get ecPrivkey failed: %v", err)
	}

	wif, err := btcutil.NewWIF(ecPrivkey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", fmt.Errorf("get wif failed: %v", err)
	}

	return wif.String(), nil
}

// normal addr version: 0x0
func GetAddress(seed []byte, nettype string, account uint32, addrIndex uint32, version uint32) (string, string, error) {

	mastkey, err := getMasterkey(seed, true) //ligo using btc mainchain cfg
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get mastkey failed: %v", err)
	}

	accountExtendKey, path, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get accountExtensionKey failed: %v", err)
	}

	accountAddr, err := getAccountAddr(accountExtendKey, nettype, version)
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get accountAddr failed: %v", err)
	}

	return accountAddr, path, nil
}

func GetAddressBytes(addr string) ([]byte, error) {
	addrBytes := base58.Decode(addr)
	if len(addrBytes) != 25 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	return addrBytes[0:21], nil
}

func GetPubkeyBytes(pub string) ([]byte, error) {
	if pub[0:1] != "1" {
		return nil, fmt.Errorf("invalid pubkey version")
	}

	pubBytes := base58.Decode(pub[1:])
	if len(pubBytes) != 37 {
		return nil, fmt.Errorf("in GetPubkeyBytes function, wrong pubkey format")
	}

	return pubBytes[:len(pubBytes)-4], nil
}

func DerivePubkey(wif string) (pub string, err error) {
	priv, err := getPrivKey(wif)
	if err != nil {
		return
	}
	buf := priv.PubKey().SerializeCompressed()
	myRipemd := ripemd160.New()

	myRipemd.Write(buf[:])
	checksum := myRipemd.Sum(nil)

	buf = append(buf, checksum[0:4]...)
	return "1" + base58.Encode(buf), nil
}

func ExportWif(seed []byte, account uint32, addrIndex uint32) (string, error) {

	mastkey, err := getMasterkey(seed, true) //ligo using btc mainchain cfg
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get mastkey failed: %v", err)
	}

	accountExtendKey, _, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get accountExtensionKey failed: %v", err)
	}

	wifKey, err := getWifkey(accountExtendKey)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get wif failed: %v", err)
	}

	return wifKey, nil

}

func getPrivKey(wif string) (*btcec.PrivateKey, error) {
	wifstruct, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, fmt.Errorf("decode wif string failed: %v", err)
	}

	return wifstruct.PrivKey, err
}

func ImportWif(wifstr string) (*btcec.PrivateKey, error) {
	return getPrivKey(wifstr)
}
