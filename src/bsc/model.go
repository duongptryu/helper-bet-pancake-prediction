package bsc

//======================= data processed

type LogProcessed struct {
	Address     string
	Bet         string
	RoundId     BigInt
	Price       float64
	BlockNumber BigInt
	TimeBefore  int64
	LogIndex BigInt
}

//==========================raw data
type ContainerLog struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []*Log  `json:"result"`
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TimeStamp        string   `json:"timeStamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}
