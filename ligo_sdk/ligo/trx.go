/**
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: trx.go
 * Date: 2018-09-04
 *
 */

package ligo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ligosdk/ligo_sdk/btssign"
)

const expireTimeout = 86000

// define trx structure
type Transaction struct {
	Ligo_ref_block_num    uint16 `json:"ref_block_num"`
	Ligo_ref_block_prefix uint32 `json:"ref_block_prefix"`
	Ligo_expiration       string `json:"expiration"`

	Ligo_operations [][]interface{} `json:"operations"`
	Ligo_extensions []interface{}   `json:"extensions"`
	Ligo_signatures []string        `json:"signatures"`

	Expiration uint32        `json:"-"`
	Operations []interface{} `json:"-"`

	Nonce uint64 `json:"nonce"`
}

func DefaultTransaction() *Transaction {

	return &Transaction{
		0,
		0,
		"",
		nil,
		nil,
		nil,
		0,
		nil,
		0,
	}
}

func GetId(id string) (uint32, error) {

	idSlice := strings.Split(id, ".")

	if len(idSlice) != 3 {
		return 0, fmt.Errorf("in GetId function, get account id failed")
	}

	res, err := strconv.ParseUint(idSlice[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("in GetId function, Parse id error %v", err)
	}

	return uint32(res), nil

}

func Str2Time(str string) int64 {

	str += "Z"
	t, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return t.Unix()

}

func Time2Str(t int64) string {

	l_time := time.Unix(t, 0).UTC()
	timestr := l_time.Format(time.RFC3339)

	timestr = timestr[:len(timestr)-1]

	return timestr
}

// in multiple precision mode
func CalculateFee(basic_op_fee int64, len_memo int64) int64 {

	var basic_memo_fee int64 = 1
	return basic_op_fee + len_memo*basic_memo_fee
}

func (asset *Asset) SetAssetBySymbol(symbol string) {
	symbol = strings.ToUpper(symbol)

	if symbol == "IGO" {
		asset.Ligo_asset_id = "1.3.0"
	} else if symbol == "BTC" {
		asset.Ligo_asset_id = "1.3.1"
	} else if symbol == "LIGO" {
		asset.Ligo_asset_id = "1.3.2"
	}

}

func GetRefblockInfo(info string) (uint16, uint32, error) {

	refinfo := strings.Split(info, ",")
	// refinfo := []string{"21771", "761216631"}

	if len(refinfo) != 2 {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, get refblockinfo failed")
	}
	ref_block_num_str, ref_block_prefix_str := refinfo[0], refinfo[1]
	ref_block_num, err := strconv.ParseUint(ref_block_num_str, 10, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_num failed: %v", err)
	}

	ref_block_prefix, err := strconv.ParseUint(ref_block_prefix_str, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_prefix failed: %v", err)
	}

	return uint16(ref_block_num), uint32(ref_block_prefix), nil
}

func GetSignature(wif string, hash []byte) ([]byte, error) {

	ecPrivkey, err := ImportWif(wif)
	if err != nil {
		return nil, fmt.Errorf("in GetSignature function, get ecprivkey failed: %v", err)
	}

	ecPrivkeyByte := ecPrivkey.Serialize()
	return btssign.SignCompact(hash, ecPrivkeyByte, true)
	//fmt.Println("the uncompressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeUncompressed()))
	//fmt.Println("the compressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeCompressed()))
	/*
		for {
			sig, err := bts.SignCompact(hash, ecPrivkeyByte, true)
			if err != nil {
				return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
			}

			pubkey_byte, err := bts.RecoverPubkey(hash, sig, true)
			if err != nil {
				return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
			}
			fmt.Println("recoverd pubkey is: ", hex.EncodeToString(pubkey_byte))

			if bytes.Compare(ecPrivkey.PubKey().SerializeCompressed(), pubkey_byte) == 0 {
				return sig, nil
			}

		}
	*/
}

func BuildUnsignedTx(refinfo, from, to, memo, assetId string, amount, fee int64, guarantee_id string) (*Transaction, error) {
	// build unsigned tx hash
	asset_amount := DefaultAsset()
	asset_amount.Ligo_amount = amount
	asset_amount.Ligo_asset_id = assetId // SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	transferOp := DefaultTransferOperation()
	transferOp.Ligo_fee = asset_fee
	transferOp.Ligo_from_addr = from
	transferOp.Ligo_to_addr = to
	transferOp.Ligo_amount = asset_amount

	if memo == "" {
		transferOp.Ligo_memo = nil
	} else {
		memo_trx := DefaultMemo()
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.Ligo_message = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))
		transferOp.Ligo_memo = &memo_trx
	}

	if guarantee_id != "" {
		transferOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return nil, err
	}

	return &Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*transferOp},
		0,
	}, nil

}

func BuildUnsignedTxHash(refinfo, from, to, memo, assetId string, amount, fee int64,
	guarantee_id, chain_id string) ([]byte, error) {
	tx, err := BuildUnsignedTx(refinfo, from, to, memo, assetId, amount, fee, guarantee_id)
	if err != nil {
		return nil, err
	}
	res := tx.Serialize()
	fmt.Printf("hex before sign: %v\n", hex.EncodeToString(res))
	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	return toSign[:], nil
}

func RebuildTxWithSign(refinfo, from, to, memo, assetId string, amount, fee int64,
	guarantee_id, sig string) ([]byte, error) {
	tx, err := BuildUnsignedTx(refinfo, from, to, memo, assetId, amount, fee, guarantee_id)
	if err != nil {
		return nil, err
	}

	tx.Ligo_signatures = append(tx.Ligo_signatures, sig)
	fmt.Printf("RebuildTxWithSign: signature=%v\n", sig)

	b, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return b, nil
}

func BuildTransferTransaction(refinfo, wif string, from, to, memo, assetId string, amount, fee int64,
	symbol string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_amount := DefaultAsset()
	asset_amount.Ligo_amount = amount
	asset_amount.Ligo_asset_id = assetId // SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	transferOp := DefaultTransferOperation()
	transferOp.Ligo_fee = asset_fee
	transferOp.Ligo_from_addr = from
	transferOp.Ligo_to_addr = to
	transferOp.Ligo_amount = asset_amount

	if memo == "" {
		transferOp.Ligo_memo = nil
	} else {
		memo_trx := DefaultMemo()
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.Ligo_message = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))
		transferOp.Ligo_memo = &memo_trx
	}

	if guarantee_id != "" {
		transferOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*transferOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// bind tunnel address fee is not needed, always 0
func BuildBindAccountTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, crosschain_wif string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	bindOp := DefaultAccountBindOperation()
	bindOp.Ligo_fee = asset_fee
	bindOp.Ligo_crosschain_type = crosschain_symbol
	bindOp.Ligo_addr = addr
	if guarantee_id != "" {
		bindOp.Ligo_guarantee_id = guarantee_id
	}

	// sign the addr
	addrByte, err := GetAddressBytes(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	toSign := sha256.Sum256(addrByte)
	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}
	bindOp.Ligo_account_signature = hex.EncodeToString(sig)
	bindOp.Ligo_tunnel_address = crosschain_addr
	crosschain_sig, err := SignAddress(crosschain_wif, crosschain_addr, crosschain_symbol)
	if err != nil {
		fmt.Println(err)
		return
	}
	bindOp.Ligo_tunnel_signature = crosschain_sig

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{10, bindOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*bindOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign = sha256.Sum256(append(chainid_byte, res...))

	sig, err = GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// BuildUnBindAccountTransaction bind tunnel address
// wif: ligo wif
// addr: ligo address
// fee:
// crosschain_addr: btc/eth/ltc/hc address
// crosschain_symbol: btc/eth/ltc/hc
// crosschain_wif: btc/eth/ltc/hc wif
// chain_id
func BuildUnBindAccountTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, crosschain_wif, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	unbindOp := DefaultAccountUnBindOperation()
	unbindOp.Ligo_fee = asset_fee
	unbindOp.Ligo_crosschain_type = crosschain_symbol
	unbindOp.Ligo_addr = addr
	//sign the addr
	addrByte, err := GetAddressBytes(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	toSign := sha256.Sum256(addrByte)
	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}
	unbindOp.Ligo_account_signature = hex.EncodeToString(sig)
	unbindOp.Ligo_tunnel_address = crosschain_addr
	crosschain_sig, err := SignAddress(crosschain_wif, crosschain_addr, crosschain_symbol)
	if err != nil {
		fmt.Println(err)
		return
	}
	unbindOp.Ligo_tunnel_signature = crosschain_sig

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{11, unbindOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*unbindOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign = sha256.Sum256(append(chainid_byte, res...))

	sig, err = GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

func BuildWithdrawCrosschainTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, assetId, crosschain_amount, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	withdrawOp := DefaultWithdrawCrosschainOperation()
	withdrawOp.Ligo_withdraw_account = addr
	withdrawOp.Ligo_amount = crosschain_amount
	withdrawOp.Ligo_asset_symbol = crosschain_symbol
	withdrawOp.Ligo_asset_id = assetId
	/*
		if crosschain_symbol == "BTC" {
			withdrawOp.Ligo_asset_symbol = "BTC"
			withdrawOp.Ligo_asset_id = "1.3.1"
		} else if crosschain_symbol == "LTC" {
			withdrawOp.Ligo_asset_symbol = "LTC"
			withdrawOp.Ligo_asset_id = "1.3.2"
		} else if crosschain_symbol == "HC" {
			withdrawOp.Ligo_asset_symbol = "HC"
			withdrawOp.Ligo_asset_id = "1.3.3"
		}
	*/
	withdrawOp.Ligo_crosschain_account = crosschain_addr

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{61, withdrawOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*withdrawOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildRegisterAccountTransaction(refinfo, wif, addr, public_key string, fee int64,
	guarantee_id, register_name, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	registerOp := DefaultRegisterAccountOperation()
	registerOp.Ligo_fee = asset_fee
	registerOp.Ligo_payer = addr
	registerOp.Ligo_name = register_name
	registerOp.Ligo_owner.Ligo_key_auths = [][]interface{}{{public_key, 1}}
	registerOp.Ligo_owner.Key_auths = public_key
	registerOp.Ligo_active.Ligo_key_auths = registerOp.Ligo_owner.Ligo_key_auths
	registerOp.Ligo_active.Key_auths = public_key
	registerOp.Ligo_owner.Key_auths = public_key
	registerOp.Ligo_options.Ligo_memo_key = public_key

	if guarantee_id != "" {
		registerOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-06T06:21:33"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{5, registerOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*registerOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildLockBalanceTransaction(refinfo, wif, addr, account_id, lock_asset_id string,
	lock_asset_amount, fee int64, miner_id, miner_address, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	lockOp := DefaultLockBalanceOperation()
	lockOp.Ligo_fee = asset_fee
	lockOp.Ligo_lock_asset_id = lock_asset_id
	lockOp.Ligo_lock_asset_amount = lock_asset_amount

	if account_id == "" {
		lockOp.Ligo_lock_balance_account = "1.2.0"
	} else {
		lockOp.Ligo_lock_balance_account = account_id
	}
	lockOp.Ligo_lock_balance_addr = addr
	lockOp.Ligo_lockto_miner_account = miner_id
	lockOp.Ligo_contract_addr = miner_address

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-07T02:18:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{55, lockOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*lockOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildRedeemBalanceTransaction(refinfo, wif, addr, account_id, foreclose_asset_id string,
	foreclose_asset_amount, fee int64, miner_id, miner_address, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	forecloseOp := DefaultForecloseBalanceOperation()
	forecloseOp.Ligo_fee = asset_fee
	forecloseOp.Ligo_foreclose_asset_id = foreclose_asset_id
	forecloseOp.Ligo_foreclose_asset_amount = foreclose_asset_amount

	forecloseOp.Ligo_foreclose_miner_account = miner_id
	forecloseOp.Ligo_foreclose_contract_addr = miner_address

	if account_id == "" {
		forecloseOp.Ligo_foreclose_account = "1.2.0"
	} else {
		forecloseOp.Ligo_foreclose_account = account_id
	}

	forecloseOp.Ligo_foreclose_addr = addr

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-07T02:18:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{56, forecloseOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*forecloseOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

// obtain_asset_arr format: []string{"citizen10,100,1.3.0", "citizen11,101,1.3.0"}
func BuildObtainPaybackTransaction(refinfo, wif, addr string, fee int64,
	obtain_asset_arr []string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Ligo_amount = fee
	asset_fee.SetAssetBySymbol("IGO")

	obtainOp := DefaultObtainPaybackOperation()
	obtainOp.Ligo_pay_back_owner = addr
	obtainOp.Ligo_fee = asset_fee

	obtainOp.Ligo_pay_back_balance = [][]interface{}{}
	if len(obtain_asset_arr) == 0 {
		return nil, fmt.Errorf("obtain asset arr forma error")
	}
	for i := 0; i < len(obtain_asset_arr); i++ {
		obtain_assets := strings.Split(obtain_asset_arr[i], ",")
		if len(obtain_assets) != 3 {
			return nil, fmt.Errorf("obtain asset arr forma error")
		}

		obtain_asset := DefaultAsset()
		amount, err := strconv.ParseInt(obtain_assets[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse obtain asset amount error")
		}
		obtain_asset.Ligo_amount = amount
		obtain_asset.Ligo_asset_id = obtain_assets[2]
		tmp_pay_back := [][]interface{}{{obtain_assets[0], obtain_asset}}
		obtainOp.Ligo_pay_back_balance = append(obtainOp.Ligo_pay_back_balance, tmp_pay_back...)
		obtainOp.citizen_name = append(obtainOp.citizen_name, obtain_assets[0])
		obtainOp.obtain_asset = append(obtainOp.obtain_asset, obtain_asset)
	}

	if guarantee_id != "" {
		obtainOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{73, obtainOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*obtainOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// fee is basic fee of ligo chain
func BuildContractInvokeTransaction(refinfo, wif, addr string, fee int64, gas_price, gas_limit int64, contract_id, contract_api, contract_arg string,
	guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	gas_count := gas_limit / 100 * gas_price
	if gas_limit%100 != 0 {
		gas_count += gas_price
	}
	asset_fee.Ligo_amount = fee + gas_count
	asset_fee.SetAssetBySymbol("IGO")

	contractOp := DefaultContractInvokeOperation()
	contractOp.Ligo_fee = asset_fee

	contractOp.Ligo_invoke_cost = uint64(gas_limit)
	contractOp.Ligo_gas_price = uint64(gas_price)
	contractOp.Ligo_caller_addr = addr
	priv, err := getPrivKey(wif)
	if err != nil {
		return nil, fmt.Errorf("get private key from wif error")
	}
	buf := priv.PubKey().SerializeCompressed()
	contractOp.Ligo_caller_pubkey = hex.EncodeToString(buf)
	contractOp.Ligo_contract_id = contract_id
	contractOp.Ligo_contract_api = contract_api
	contractOp.Ligo_contract_arg = contract_arg

	if guarantee_id != "" {
		contractOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{79, contractOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*contractOp},
		0,
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// transfer to contract
func BuildContractTransferTransaction(refinfo, wif, addr string, fee int64, amount int64, assetId string, gas_price, gas_limit int64, contract_id, param string,
	guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Ligo_amount = CalculateFee(2000, int64(len(memo) + 3))
	gas_count := gas_limit / 100 * gas_price
	if gas_limit%100 != 0 {
		gas_count += gas_price
	}
	asset_fee.Ligo_amount = fee + gas_count
	asset_fee.SetAssetBySymbol("IGO")

	asset_amount := DefaultAsset()
	asset_amount.Ligo_amount = amount
	asset_amount.Ligo_asset_id = assetId

	contractOp := DefaultContractTransferOperation()
	contractOp.Ligo_fee = asset_fee
	contractOp.Ligo_amount = asset_amount

	contractOp.Ligo_invoke_cost = uint64(gas_limit)
	contractOp.Ligo_gas_price = uint64(gas_price)
	contractOp.Ligo_caller_addr = addr
	priv, err := getPrivKey(wif)
	if err != nil {
		return nil, fmt.Errorf("get private key from wif error")
	}
	buf := priv.PubKey().SerializeCompressed()
	contractOp.Ligo_caller_pubkey = hex.EncodeToString(buf)
	contractOp.Ligo_contract_id = contract_id
	contractOp.Ligo_param = param

	if guarantee_id != "" {
		contractOp.Ligo_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{81, contractOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*contractOp},
		0,
	}

	res := transferTrx.Serialize()

	fmt.Println("chain_id:", chain_id)
	fmt.Println("res:", hex.EncodeToString(res))

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))
	fmt.Println("wif:", wif, "toSign:", hex.EncodeToString(toSign[:]))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Ligo_signatures = append(transferTrx.Ligo_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("tx:", string(b))
	return
}
