package bsc

import (
	"by_go/src/hexutil"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type client struct {
	host string
	apiKey string

	locker *sync.RWMutex
	cache map[string]string
}

func NewClient (apiKey string) *client {
	return &client{
		host: HostApiBsc,
		apiKey: apiKey,
		cache: make(map[string]string),
		locker: new(sync.RWMutex),
	}
}

func (c *client) GetLogs (fromBlock string, toBlockFake *string, address string) ([]*LogProcessed, error) {
	var toBlock string

	if toBlockFake == nil {
		toBlock = "latest"
	}else {
		toBlock = *toBlockFake
	}

	params := fmt.Sprintf("?module=logs&action=getLogs&fromBlock=%s&toBlock=%s&address=%s&apikey=%s", fromBlock, toBlock, address, c.apiKey)

	url := c.host + params
	result, err := makeGetLogHttpRequest(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var logProcessed []*LogProcessed
	for _, data := range result {
		if _, ok := c.cache[data.TransactionHash]; ok {
			continue
		}else {
			c.locker.RLock()
			c.cache[data.TransactionHash] = data.TransactionIndex
			c.locker.RUnlock()
		}

		var l LogProcessed
		if data.Topics[0] == BearDown {
			l.Bet = "Down"
		}else if data.Topics[0] == BullUp {
			l.Bet = "Up"
		}else {
			continue
		}

		timeNow := time.Now().Unix()
		timeBet, _ := hexutil.DecodeBig(data.TimeStamp)
		l.TimeBefore = timeNow - timeBet.Int64()

		i := new(big.Float)
		err := i.UnmarshalText([]byte(data.Data))
		if err != nil {
			log.Fatalln(err)
		}
		n, _ := i.Float64()
		l.Price = n/ DeviceNumber


		roundId := new(BigInt)
		err = roundId.UnmarshalText([]byte(data.Topics[2]))
		if err != nil {
			log.Fatalln(err, data.Topics[2])
		}
		l.RoundId = *roundId

		blockNumber := new(BigInt)
		_ = blockNumber.UnmarshalText([]byte(data.BlockNumber))
		l.BlockNumber = *blockNumber

		l.Address = "https://bscscan.com/address/" + data.Topics[1]

		logIndex := new(BigInt)
		_ = logIndex.UnmarshalText([]byte(data.LogIndex))
		l.LogIndex = *logIndex

		logProcessed = append(logProcessed, &l)
	}

	return logProcessed, nil
}

func makeGetLogHttpRequest(url string) ([]*Log, error) {
	httpClient := http.Client{
		Timeout: 5*time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status request != 200 -%d", resp.StatusCode))
	}

	var result ContainerLog
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}
