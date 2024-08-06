package ligo_sdk

import (
	"fmt"
	"ligosdk/ligo_sdk/ligo"
	"testing"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
)

var (
	//ligouri      = "http://192.168.1.220/server/process"
	ligouri   = "http://192.168.1.220/server/process"
	walletURI = "http://192.168.1.122:10046" // broadcast wallet rpc api
	// zhengqinpeng wid
	ligowalletid = "dcc02a08142e230814e108053eb40c6180f908da"
	//myligoaddr   = "LIGONhNthqhgkEfPjzhLQ3cWnNEpEjpjtvmKzw"
	//myhcaddr   = "TsTaA9DT2z24Wg4U31QLeSvcg1agu6e4d88"
	// zhengqinpeng address
	myligoaddr = "" // "LIGONWjKv1PUbZ6dgVoerkohdoCHRck9LAZh3Y"
	// myligoaddr = "LIGONXPFrzwsTo7wT5QceDLEjpd2VZ8BZR3UCY"  // just for bind usdt
	// wifligo = "5KZTntQJ9AAARnuKUzgUG4MHiPPVU7FuMq4d3WNH4rhasDxRa4a" // for bind usdt

	myltcaddr   = "LP2JKjy9WmSygMdoe2CzEHabXrPSPXspNF" // "LQz1gokPokkj6hK3dHt1JphZbPh5G9KpbV"
	dstligoaddr = "LIGONL3kJ4prkHUsHsnwW4HGSDUfYxncWcfgDn"

	mycointype = "LIGO"                                //  "HC(LIGO)"
	myhcaddr   = "TsZj5Cx1p94izYoxVuZyRGvskSdGaasCHir" // hc testnet
	myxwif     = "PtWUvDcPuchqQA5vfoUqoNJw6JyCSnYrbHM37XXVUbfLT4wdJj7hU"

	//mycointype = "ETH(LIGO)"
	//myhcaddr  = "0x1891025831596418915523e786334b2b44985272"
	//myxwif = "3f2153c638e857ae4b5ef132c1ee09c24bb48484d2dea91a5071b202be2e2a90" // eth

	//mycointype = "PAX(LIGO)"
	//myhcaddr  = "0x1891025831596418915523e786334b2b44985272"
	//myxwif = "3f2153c638e857ae4b5ef132c1ee09c24bb48484d2dea91a5071b202be2e2a90" // eth

	widbtc       = "a8eb57aacdb63a6ed8d485c0304260d7e627d704"
	cointypebtc  = "BTC(LIGO)"
	coinaddrbtc  = "1PNvmFKPGDADmrPXQcLVhyqFPvSZ1czHa"
	wifbtc       = "Kx9XqhpA2UTJm21HgVc3yuvcGMQNfyKi133MEUNGhb6NvNniLzH8" // btc
	cointypeusdt = "USDT(LIGO)"
	coinaddrusdt = "1PNvmFKPGDADmrPXQcLVhyqFPvSZ1czHa"
	wifusdt      = "L5QrMWSdpU8CQer7ss5zofTzmzgcB7spmBSSoXF29NKNZVYKxwkB"

	//mycointype = "LTC(LIGO)"
	//myhcaddr  = "LiSt5dsfy7WB5Yq1AUaxjjHPpNQU8upKES"
	//myxwif = "6uHv9Ru7dFrv1easCC1sRwvFSf9gFm2FUpBC1P2TCT9pVQyUZCo" // ltc

	//mycointype = "LTC(LIGO)"
	//myhcaddr  = "LP2JKjy9WmSygMdoe2CzEHabXrPSPXspNF"
	//myxwif = "T5LN3q4o3CUtKpoi3TtiL4JMPNjqaBhkoWzQJ2F4HLWuu94Y6G8y" // ltc
	ligo_mne     = "hobby actual sadness know copy achieve bulb message unhappy snack giggle core reason enroll boat magic aim sea front capital text science green joy"
	ligo_net     = "LIGO"
	ligo_chainId = "5200ea0fc76d785ec205805fd287d3b28cea78f4db58fe41cd833077f20b0ffb"
)

const (
	CoinLIGO string = "LIGO"

	VersionNormalAddr   = 0x35
	VersionMultisigAddr = 0x32
	VersionContractAddr = 0x1c
)

func init() {
	seed := ligo.MnemonicToSeed(ligo_mne, "")
	addr, _, err := ligo.GetAddress([]byte(seed), "LIGO", 0, 0, 0)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("derived ligo address is %v\n", addr)
	myligoaddr = addr
}

func TestLIGOAddress(t *testing.T) {

	for i := 0; i < 5; i++ {
		seed := ligo.MnemonicToSeed(ligo_mne, "")
		addr, _, err := ligo.GetAddress([]byte(seed), "LIGO", 0, uint32(i), 0)
		assert.Nil(t, err)
		fmt.Printf("LIGO index %d Address: %v\n", i, addr)
	}
}

func TestLIGOPubkey(t *testing.T) {

	seed := ligo.MnemonicToSeed(ligo_mne, "")
	wif, err := ligo.ExportWif(seed, uint32(0), uint32(0))
	assert.Nil(t, err)

	pk, err := ligo.DerivePubkey(wif)
	assert.Nil(t, err)
	fmt.Println(pk)

	// test
	wif = "5KR6ocp5eUdWWYPX7mYp4XLGBcZ2xHVHVsNaco6K2YZSWQTqES7"
	pk, err = ligo.DerivePubkey(wif)
	assert.Nil(t, err)
	fmt.Println(pk)
	assert.Equal(t, "15cj54KW1TDK94GCcrSovYDUUYQ1FgCbECnG2CaEv5ub7Vghx5N", pk)
}

func _doReq(t *testing.T, uri string, ms map[string]interface{}) []byte {
	body, err := json.Marshal(ms)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	assert.Nil(t, err)

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true
	// req.

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, 200)

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	trancode := ""
	iheader := ms["header"]
	if iheader != nil {
		if header, ok := iheader.(map[string]interface{}); ok {
			itrancode, ok := header["trancode"]
			if ok {
				trancode = itrancode.(string)
			}
		}
	}
	fmt.Printf("request %s %s: req body=%v\ncode=%d response=%v\n", uri, trancode, string(body), res.StatusCode, string(buf))
	assert.Nil(t, err)

	return buf
}

func doReq(t *testing.T, ms map[string]interface{}) []byte {
	return _doReq(t, ligouri, ms)
}

func makeheader(transcode, walletid string) map[string]interface{} {
	m := map[string]interface{}{}
	m["header"] = map[string]string{
		"version":    "2.4.0",
		"language":   "zh-Hans",
		"trancode":   transcode,
		"clienttype": "Android",
		"walletid":   walletid,
		"random":     "abc",
		"handshake":  "efg",
		"imie":       "3511",
		"source":     "sdk",
	}
	m["body"] = map[string]interface{}{}

	return m
}

func setBodyData(t *testing.T, dst map[string]interface{}, buf []byte) {
	var m map[string]interface{}

	err := json.Unmarshal(buf, &m)
	assert.Nil(t, err)

	dst["data"] = m
}

func directPostWallet(t *testing.T, action string, buf []byte) {
	_directPostWallet(t, walletURI, buf)
}

func _directPostWallet(t *testing.T, uri string, buf []byte) []byte {
	// uri := walletURI // "http://192.168.1.128:10033"
	param := map[string]interface{}{
		"id":     1,
		"method": "lightwallet_broadcast",
	}
	var m map[string]interface{}
	err := json.Unmarshal(buf, &m)
	assert.Nil(t, err)

	param["params"] = []interface{}{m}
	res := _doReq(t, uri, param)
	return res
	// fmt.Printf("action %v response: %v\n", action, string(res))
}

// test bind BTC
// test bind USDT, address is same with BTC
