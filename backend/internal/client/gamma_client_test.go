package client

import (
	"encoding/json"
	"testing"
)

func TestDecodeGammaMarketsVariants(t *testing.T) {
	base := []GammaMarket{{ConditionID: "0x1", Question: "q"}}
	cases := []json.RawMessage{
		json.RawMessage(`[{"conditionId":"0x1","question":"q"}]`),
		json.RawMessage(`{"data":[{"conditionId":"0x1","question":"q"}]}`),
		json.RawMessage(`{"markets":[{"conditionId":"0x1","question":"q"}]}`),
		json.RawMessage(`[{"conditionId":"0x1","question":"q","volume":"32257.445115","liquidity":"12.5","active":true}]`),
		json.RawMessage(`[{"conditionId":"0x1","question":"q","volume":32257.445115,"liquidity":12.5,"active":true}]`),
	}

	for i, c := range cases {
		got, err := decodeGammaMarkets(c)
		if err != nil {
			t.Fatalf("case %d decode err: %v", i, err)
		}
		if len(got) != len(base) || got[0].ConditionID != base[0].ConditionID {
			t.Fatalf("case %d unexpected decode: %#v", i, got)
		}
	}
}

func TestDecodeGammaMarketsParsesStringNumbers(t *testing.T) {
	raw := json.RawMessage(`[{"conditionId":"0x1","question":"q","volume":"10.25","liquidity":"3.5","active":true}]`)
	got, err := decodeGammaMarkets(raw)
	if err != nil {
		t.Fatalf("decode err: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("unexpected len: %d", len(got))
	}
	if got[0].Volume != 10.25 || got[0].Liquidity != 3.5 {
		t.Fatalf("unexpected numeric decode: %#v", got[0])
	}
}
