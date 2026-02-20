package client

import (
	"encoding/json"
	"testing"
)

func TestDecodeTradesVariants(t *testing.T) {
	cases := []json.RawMessage{
		json.RawMessage(`[{"transactionHash":"0x1","timestamp":1,"market":"m"}]`),
		json.RawMessage(`{"data":[{"transactionHash":"0x1","timestamp":1,"market":"m"}]}`),
		json.RawMessage(`{"trades":[{"transactionHash":"0x1","timestamp":1,"market":"m"}]}`),
	}

	for i, c := range cases {
		got, err := decodeTrades(c)
		if err != nil {
			t.Fatalf("case %d decode err: %v", i, err)
		}
		if len(got) != 1 || got[0].TransactionHash != "0x1" {
			t.Fatalf("case %d unexpected decode: %#v", i, got)
		}
	}
}
