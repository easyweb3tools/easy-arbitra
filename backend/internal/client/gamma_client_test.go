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
