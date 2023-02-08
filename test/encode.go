package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func main() {
	addr := "TYtTjZCvEMKQ42dxeYSyDsWx2SkEQFwqPh"
	decode, v, err := base58.CheckDecode(addr)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}

	hexTo := fmt.Sprintf("41%x", decode)
	fmt.Printf("dst:%s,%x\n", hexTo, v)

	//反向
	//hex, err := hex.DecodeString(hexTo[2:])
	hex, err := hex.DecodeString("7f201283bb4c475f1dee17847f12d27445008132")
	src := base58.CheckEncode(hex, 65)
	fmt.Printf("src:%s\n", src)
}
