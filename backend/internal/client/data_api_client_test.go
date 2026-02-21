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
		json.RawMessage(`[{"transactionHash":"0x1","timestamp":1,"conditionId":"0xc1","asset":"a1","proxyWallet":"0xabc"}]`),
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

func TestDecodeTradesUsesConditionAndProxyWallet(t *testing.T) {
	raw := json.RawMessage(`[{"transactionHash":"0x1","timestamp":1,"conditionId":"0xc1","asset":"a1","proxyWallet":"0xabc"}]`)
	got, err := decodeTrades(raw)
	if err != nil {
		t.Fatalf("decode err: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("unexpected len: %d", len(got))
	}
	if got[0].Market != "0xc1" {
		t.Fatalf("expected market from conditionId, got: %s", got[0].Market)
	}
	if got[0].TakerAddress != "0xabc" {
		t.Fatalf("expected taker from proxyWallet, got: %s", got[0].TakerAddress)
	}
}
