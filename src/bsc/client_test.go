package bsc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClient_GetLogs(t *testing.T) {
	c := NewClient("FCAHNM2UJIWS4ADYQM9U1D1VCTFFPIFV55")
	fromBlock := "14182279"
	data, _ := c.GetLogs(fromBlock, nil, "0x18B2A687610328590Bc8F2e5fEdDe3b582A49cdA")
	dump(data)
}

func dump(v interface{}) {
	raw, err := json.MarshalIndent(v, "", "	")
	if nil != err {
		fmt.Println("Dump error: ", err)
	} else {
		fmt.Println(string(raw))
	}
}