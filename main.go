package main

import (
	"fmt"
	"ligosdk/ligo_sdk"
)

// Kxw17Y8T11kNrbaY8Y53aXkNvRo8tgYJGZaAYf9bUDBQKkfXXM3z LIGO5TDS4UrrUTAjmz5sQafYUrM37obZvCEyrVxJHd6teq5wiB7UDA LIGONTyhBEVF312RfTyoQ878AhQwerayc7eazr <nil>
func main() {
	fmt.Println("just testing")
	//wif, pubkey, addr, error := ligo_sdk.GetNewPrivate()
	//fmt.Println(wif, pubkey, addr, error)
	wif := "L5cA1ui3UToWJdVCdVKqppNYxFoWV21JpCJyBm9JRdWEuCf8jV6x"
	//pubkey := "18RyCsazpmJ3xGHzzUVGdSjR3UjF7EqB8Kw9JrCg74tM4TrFdrt"
	addr := "13kYzxE2BEizhBr8Z2orM5xpywjoz8TSLm"
	to_addr := "15iq2e57JV3yURHQiTPxNV72Wuekorvexa"
	ref_info := ligo_sdk.CalRefInfo("001c7fe6483a14fc3d553aa23140eeb853874c59")
	fmt.Println(ref_info)
	ligo_chain_id := "5200ea0fc76d785ec205805fd287d3b28cea78f4db58fe41cd833077f20b0ffb"
	trx_data, err := ligo_sdk.LigoTransfer(ref_info, wif, ligo_chain_id, addr, to_addr, "IGO", "0.11", "1", "", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("")
	fmt.Println("raw trx", string(trx_data))
	trx_data, err = ligo_sdk.LigoRegister(ref_info, wif, ligo_chain_id, "newtest", addr, "5000", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("")
	fmt.Println("register raw trx", string(trx_data))
	trx_data, err = ligo_sdk.LigoMining(ref_info, wif, ligo_chain_id, "IGO", "1.2.105", addr, "1", "0", "1.6.1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("")
	fmt.Println("LigoMining raw trx", string(trx_data))
	trx_data, err = ligo_sdk.LigoForecloseBalance(ref_info, wif, ligo_chain_id, addr, "1.2.105", "1.3.0", "1.6.1", "1", "0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("")
	fmt.Println("LigoForecloseBalance raw trx", string(trx_data))

}
