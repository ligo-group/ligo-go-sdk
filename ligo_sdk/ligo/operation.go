/**
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: operation.go.go, Date: 2018-10-31
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package ligo

type Asset struct {
	Ligo_amount   int64  `json:"amount"`
	Ligo_asset_id string `json:"asset_id"`
}

// ligo  --- "1.3.0"
// btc --- "1.3.1"
// ltc --- "1.3.2"
// hc  --- "1.3.3"
func DefaultAsset() Asset {
	return Asset{
		0,
		"1.3.0",
	}
}

type Extension struct {
	extension []string
}

type Memo struct {
	Ligo_from    string `json:"from"` //public_key_type  33
	Ligo_to      string `json:"to"`   //public_key_type  33
	Ligo_nonce   uint64 `json:"nonce"`
	Ligo_message string `json:"message"`

	IsEmpty bool   `json:"-"`
	Message string `json:"-"`
}

func DefaultMemo() Memo {

	return Memo{
		"LIGO1111111111111111111111111111111114T1Anm",
		"LIGO1111111111111111111111111111111114T1Anm",
		0,
		"",
		true,
		"",
	}

}

type Authority struct {
	Ligo_weight_threshold uint32          `json:"weight_threshold"`
	Ligo_account_auths    []interface{}   `json:"account_auths"`
	Ligo_key_auths        [][]interface{} `json:"key_auths"`
	Ligo_address_auths    []interface{}   `json:"address_auths"`

	Key_auths string `json:"-"`
}

func DefaultAuthority() Authority {

	return Authority{
		1,
		[]interface{}{},
		[][]interface{}{{"", 1}},
		[]interface{}{},
		"",
	}
}

type AccountOptions struct {
	Ligo_memo_key              string        `json:"memo_key"`
	Ligo_voting_account        string        `json:"voting_account"`
	Ligo_num_witness           uint16        `json:"num_witness"`
	Ligo_num_committee         uint16        `json:"num_committee"`
	Ligo_votes                 []interface{} `json:"votes"`
	Ligo_miner_pledge_pay_back byte          `json:"miner_pledge_pay_back"`
	Ligo_extensions            []interface{} `json:"extensions"`
}

func DefaultAccountOptions() AccountOptions {

	return AccountOptions{
		"",
		"1.2.5",
		0,
		0,
		[]interface{}{},
		10,
		[]interface{}{},
	}

}

// transfer operation tag is  0
type TransferOperation struct {
	Ligo_fee          Asset  `json:"fee"`
	Ligo_guarantee_id string `json:"guarantee_id,omitempty"`
	Ligo_from         string `json:"from"`
	Ligo_to           string `json:"to"`

	Ligo_from_addr string `json:"from_addr"`
	Ligo_to_addr   string `json:"to_addr"`

	Ligo_amount Asset `json:"amount"`
	Ligo_memo   *Memo `json:"memo,omitempty"`

	Ligo_extensions []interface{} `json:"extensions"`
}

func DefaultTransferOperation() *TransferOperation {

	return &TransferOperation{
		DefaultAsset(),
		"",
		"1.2.0",
		"1.2.0",
		"",
		"",
		DefaultAsset(),
		nil,
		make([]interface{}, 0),
	}
}

// account bind operation tag is 10
type AccountBindOperation struct {
	Ligo_fee               Asset  `json:"fee"`
	Ligo_crosschain_type   string `json:"crosschain_type"`
	Ligo_addr              string `json:"addr"`
	Ligo_account_signature string `json:"account_signature"`
	Ligo_tunnel_address    string `json:"tunnel_address"`
	Ligo_tunnel_signature  string `json:"tunnel_signature"`
	Ligo_guarantee_id      string `json:"guarantee_id,omitempty"`
}

func DefaultAccountBindOperation() *AccountBindOperation {

	return &AccountBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// account unbind operation tag is 11
type AccountUnBindOperation struct {
	Ligo_fee               Asset  `json:"fee"`
	Ligo_crosschain_type   string `json:"crosschain_type"`
	Ligo_addr              string `json:"addr"`
	Ligo_account_signature string `json:"account_signature"`
	Ligo_tunnel_address    string `json:"tunnel_address"`
	Ligo_tunnel_signature  string `json:"tunnel_signature"`
}

func DefaultAccountUnBindOperation() *AccountUnBindOperation {

	return &AccountUnBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
	}
}

// withdraw cross chain operation tag is 61
type WithdrawCrosschainOperation struct {
	Ligo_fee              Asset  `json:"fee"`
	Ligo_withdraw_account string `json:"withdraw_account"`
	Ligo_amount           string `json:"amount"`
	Ligo_asset_symbol     string `json:"asset_symbol"`

	Ligo_asset_id           string `json:"asset_id"`
	Ligo_crosschain_account string `json:"crosschain_account"`
	Ligo_memo               string `json:"memo"`
}

func DefaultWithdrawCrosschainOperation() *WithdrawCrosschainOperation {

	return &WithdrawCrosschainOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// register account operation tag is 5
type RegisterAccountOperation struct {
	Ligo_fee              Asset     `json:"fee"`
	Ligo_registrar        string    `json:"registrar"`
	Ligo_referrer         string    `json:"referrer"`
	Ligo_referrer_percent uint16    `json:"referrer_percent"`
	Ligo_name             string    `json:"name"`
	Ligo_owner            Authority `json:"owner"`
	Ligo_active           Authority `json:"active"`
	Ligo_payer            string    `json:"payer"`

	Ligo_options      AccountOptions `json:"options"`
	Ligo_extensions   interface{}    `json:"extensions"`
	Ligo_guarantee_id string         `json:"guarantee_id,omitempty"`
}

func DefaultRegisterAccountOperation() *RegisterAccountOperation {

	return &RegisterAccountOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",
		0,
		"",
		DefaultAuthority(),
		DefaultAuthority(),
		"",

		DefaultAccountOptions(),
		make(map[string]interface{}, 0),
		"",
	}

}

// lock balance operation tag is 55
type LockBalanceOperation struct {
	Ligo_lock_asset_id     string `json:"lock_asset_id"`
	Ligo_lock_asset_amount int64  `json:"lock_asset_amount"`
	Ligo_contract_addr     string `json:"contract_addr"`

	Ligo_lock_balance_account string `json:"lock_balance_account"`
	Ligo_lockto_miner_account string `json:"lockto_miner_account"`
	Ligo_lock_balance_addr    string `json:"lock_balance_addr"`

	Ligo_fee Asset `json:"fee"`
}

func DefaultLockBalanceOperation() *LockBalanceOperation {

	return &LockBalanceOperation{
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
		DefaultAsset(),
	}
}

// foreclose balance operation tag is 56
type ForecloseBalanceOperation struct {
	Ligo_fee Asset `json:"fee"`

	Ligo_foreclose_asset_id     string `json:"foreclose_asset_id"`
	Ligo_foreclose_asset_amount int64  `json:"foreclose_asset_amount"`

	Ligo_foreclose_miner_account string `json:"foreclose_miner_account"`
	Ligo_foreclose_contract_addr string `json:"foreclose_contract_addr"`

	Ligo_foreclose_account string `json:"foreclose_account"`
	Ligo_foreclose_addr    string `json:"foreclose_addr"`
}

func DefaultForecloseBalanceOperation() *ForecloseBalanceOperation {

	return &ForecloseBalanceOperation{
		DefaultAsset(),
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
	}
}

// obtain pay back operation tag is 73
type ObtainPaybackOperation struct {
	Ligo_pay_back_owner   string          `json:"pay_back_owner"`
	Ligo_pay_back_balance [][]interface{} `json:"pay_back_balance"`
	Ligo_guarantee_id     string          `json:"guarantee_id,omitempty"`
	Ligo_fee              Asset           `json:"fee"`

	citizen_name []string
	obtain_asset []Asset
}

func DefaultObtainPaybackOperation() *ObtainPaybackOperation {

	return &ObtainPaybackOperation{
		"",
		[][]interface{}{{"", DefaultAsset()}},
		"",
		DefaultAsset(),
		nil,
		nil,
	}
}

// contract invoke operation tag is 79
type ContractInvokeOperation struct {
	Ligo_fee           Asset  `json:"fee"`
	Ligo_invoke_cost   uint64 `json:"invoke_cost"`
	Ligo_gas_price     uint64 `json:"gas_price"`
	Ligo_caller_addr   string `json:"caller_addr"`
	Ligo_caller_pubkey string `json:"caller_pubkey"`
	Ligo_contract_id   string `json:"contract_id"`
	Ligo_contract_api  string `json:"contract_api"`
	Ligo_contract_arg  string `json:"contract_arg"`
	//Ligo_extension     []interface{} `json:"extensions"`
	Ligo_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractInvokeOperation() *ContractInvokeOperation {

	return &ContractInvokeOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		"",
		"",
		//make([]interface{}, 0),
		"",
	}
}

// transfer to contract operation tag is 81
type ContractTransferOperation struct {
	Ligo_fee           Asset  `json:"fee"`
	Ligo_invoke_cost   uint64 `json:"invoke_cost"`
	Ligo_gas_price     uint64 `json:"gas_price"`
	Ligo_caller_addr   string `json:"caller_addr"`
	Ligo_caller_pubkey string `json:"caller_pubkey"`
	Ligo_contract_id   string `json:"contract_id"`
	Ligo_amount        Asset  `json:"amount"`
	Ligo_param         string `json:"param"`
	//Ligo_extension     []interface{} `json:"extensions"`
	Ligo_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractTransferOperation() *ContractTransferOperation {

	return &ContractTransferOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		DefaultAsset(),
		"",
		//make([]interface{}, 0),
		"",
	}
}
