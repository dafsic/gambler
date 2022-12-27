package modules

type Tx struct {
	Hash      string `json:"transactionId"`
	Height    int64  `json:"blockNumber"`
	BlockHash string `json:"blockHash"`
	Result    string `json:"result"`
}

type Block struct {
	Ts        int64    `json:"timeStamp"`
	Height    int64    `json:"blockNumber"`
	BlockHash string   `json:"blockHash"`
	Txs       []string `json:"transactionList"`
}
