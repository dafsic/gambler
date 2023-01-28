package main

import (
	"fmt"

	"github.com/dafsic/gambler/utils"
)

var key = "iDxSd9m8zJ6wyCh7"
var pdata = `{"From":"41ee7f5f79e8cac83c7942120c088630ec5ef47fa0","To":"417e36ce97d6f8a5e47b1c95ae01fe815bc8f2c8cd","Token":"TRX","Amount":"1000000","Ts":1672888075}`

func main() {
	//加密
	c, e := utils.AesEncryptoBase64(pdata, key)
	fmt.Printf("加密:%s,%v\n", c, e)

	//解密
	d, e := utils.AesDecryptooBase64(c, key)
	fmt.Printf("解密:%s,%v\n", d, e)
}
