package robot

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/utils"
)

var oddOrEvenMap = map[byte]bool{'0': true, '1': false, '2': true, '3': false, '4': true, '5': false, '6': true, '7': false, '8': true, '9': false}

// IsEvenNum 判断一个字符串最后一个数字是偶数吗
func IsEvenNum(h string) bool {
	l := len(h)
	for i := 1; i <= l; i++ {
		if h[l-i] > '9' {
			continue
		}
		return oddOrEvenMap[h[l-i]]
	}
	return true
}

func IsRefund(from, to string, minTs, maxTs int64, token string) (bool, error) {
	switch strings.ToUpper(token) {
	case "TRX":
		return IsRefundTrx(from, to, minTs, maxTs)
	case "USDT":
		return IsRefundUsdt(from, to, minTs, maxTs, modules.USDT_CONTRACT)
	default:
		return false, nil
	}
}

func IsRefundTrx(from, to string, minTs, maxTs int64) (bool, error) {
	//url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions?only_to=true&min_timestamp=%d&max_timestamp=%d&search_internal=false", to, minTs, maxTs)
	url := fmt.Sprintf("https://nile.trongrid.io/v1/accounts/%s/transactions?only_to=true&min_timestamp=%d&max_timestamp=%d&search_internal=false", to, minTs, maxTs)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return strings.Contains(string(body), from), nil
}
func IsRefundUsdt(from, to string, minTs, maxTs int64, ca string) (bool, error) {
	//url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?only_to=true&min_timestamp=%d&max_timestamp=%d&search_internal=false&contract_address=%s", to, minTs, maxTs)
	url := fmt.Sprintf("https://nile.trongrid.io/v1/accounts/%s/transactions/trc20?only_to=true&min_timestamp=%d&max_timestamp=%d&search_internal=false&contract_address=%s", to, minTs, maxTs, hex2addr(ca))

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return strings.Contains(string(body), hex2addr(from)), nil
}

func hex2addr(h string) string {
	a, err := hex.DecodeString(h[2:])
	if err != nil {
		return ""
	}

	return base58.CheckEncode(a, 65)
}

func addr2hex(a string) string {
	decode, _, err := base58.CheckDecode(a)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("41%x", decode)
}
