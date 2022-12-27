package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var data = []byte(`{"data":[{"ret":[{"contractRet":"SUCCESS","fee":268000}],"signature":["919888f663066f69000f7e714d95ba077c52c69c3ce1645a48cf02adc5bfde6d4258abb34a671e35ec17d5f6c9f80dfc5770ad0c1e1e3ee666d4c645d230867901"],"txID":"92bd39bf617c01867f0404aa170eb463b9a13b715fa96afe0e99d5dec7c06c8d","net_usage":0,"raw_data_hex":"0a024eee2208a9e5c8f8958f16614080f8e3e4d5305a68080112640a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412330a1541de08bafb1cdc90280b4687c0f21e55e23a7764c31215414b838016627eec39be80e570cd2db4249320aa6b18cbcd8c0370d7d7c4bbd530","net_fee":268000,"energy_usage":0,"blockNumber":47206147,"block_timestamp":1672209471000,"energy_fee":0,"energy_usage_total":0,"raw_data":{"contract":[{"parameter":{"value":{"amount":6497995,"owner_address":"41de08bafb1cdc90280b4687c0f21e55e23a7764c3","to_address":"414b838016627eec39be80e570cd2db4249320aa6b"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"4eee","ref_block_hash":"a9e5c8f8958f1661","expiration":1672295808000,"timestamp":1672209312727},"internal_transactions":[]},{"ret":[{"contractRet":"SUCCESS","fee":0}],"signature":["d91a128505a4cbecb3c06f501e01f693ea328b85113cb5fb3e7e990e0266e4f92789b756636934668d0915bcb5c434c36d2407289c4c2ebe784d58df0c22b02200"],"txID":"8c0f529cc2d1417e18fa2c85ef081701ebcc0ff45d8162883ae87a342be9a17d","net_usage":280,"raw_data_hex":"0a024edb2208f8da23af63732cf840a8fdcdbbd5305a74080212700a32747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e736665724173736574436f6e7472616374123a0a073130303439303512154184bfd8a960edd27f6bfeb358eb26f654d1c66eba1a15414b838016627eec39be80e570cd2db4249320aa6b20f8067095b9cabbd530","net_fee":0,"energy_usage":0,"blockNumber":47206126,"block_timestamp":1672209408000,"energy_fee":0,"energy_usage_total":0,"raw_data":{"contract":[{"parameter":{"value":{"amount":888,"asset_name":"1004905","owner_address":"4184bfd8a960edd27f6bfeb358eb26f654d1c66eba","to_address":"414b838016627eec39be80e570cd2db4249320aa6b"},"type_url":"type.googleapis.com/protocol.TransferAssetContract"},"type":"TransferAssetContract"}],"ref_block_bytes":"4edb","ref_block_hash":"f8da23af63732cf8","expiration":1672209465000,"timestamp":1672209407125},"internal_transactions":[]},{"ret":[{"contractRet":"SUCCESS","fee":0}],"signature":["4578463879c048daa8cc90a6935daf4c0dd1ca119d337b69fdb933f949de04d2335b39ae47881265005ea1241d25bdcff39d608cb7c46e8279ac89a6c5edc76e00"],"txID":"ee21169461206cb40c17754df7efbab99be5a64b4e960d6eb077e77a743452fd","net_usage":265,"raw_data_hex":"0a024edb2208f8da23af63732cf840a8fdcdbbd5305a65080112610a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412300a1541f6744faebf55da0e8e66c197c88f98cd32cec5371215414b838016627eec39be80e570cd2db4249320aa6b1801708eb0cabbd530","net_fee":0,"energy_usage":0,"blockNumber":47206126,"block_timestamp":1672209408000,"energy_fee":0,"energy_usage_total":0,"raw_data":{"contract":[{"parameter":{"value":{"amount":1,"owner_address":"41f6744faebf55da0e8e66c197c88f98cd32cec537","to_address":"414b838016627eec39be80e570cd2db4249320aa6b"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"4edb","ref_block_hash":"f8da23af63732cf8","expiration":1672209465000,"timestamp":1672209405966},"internal_transactions":[]}],"success":true,"meta":{"at":1672289611202,"page_size":3}}`)

	var r Result

	e := json.Unmarshal(data, &r)
	if e != nil {
		fmt.Println("1111:", e.Error())
		return
	}

	fmt.Println("2222:", r.D[0])

}

type Result struct {
	D []Data `json:"data"`
}

type Data struct {
	R RawData `json:"raw_data"`
}

type RawData struct {
	C []Contract `json:"contract"`
}

type Contract struct {
	P Parameter `json:"parameter"`
}

type Parameter struct {
	V Value `json:"value"`
}

type Value struct {
	Owner  string `json:"owner_address"`
	Amount int64  `json:"amount"`
}
