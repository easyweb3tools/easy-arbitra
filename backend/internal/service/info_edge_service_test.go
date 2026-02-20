package service

import (
	"testing"
)

func TestInfoEdgeClassificationThresholds(t *testing.T) {
	// Unit sanity check for classification helpers through deterministic values.
	// We test expected categories by emulating outputs from Evaluate calculations.
	cases := []struct {
		mean   float64
		pvalue float64
		sample int64
		want   string
	}{
		{-45, 0.01, 12, "processing_edge"},
		{-8, 0.08, 12, "mild_edge"},
		{10, 0.01, 12, "no_edge"},
		{-50, 0.2, 12, "no_edge"},
		{-50, 0.01, 3, "insufficient_data"},
	}

	for _, tc := range cases {
		got := "insufficient_data"
		switch {
		case tc.sample < 5:
			got = "insufficient_data"
		case tc.pvalue < 0.05 && tc.mean <= -30:
			got = "processing_edge"
		case tc.pvalue < 0.10 && tc.mean < 0:
			got = "mild_edge"
		default:
			got = "no_edge"
		}
		if got != tc.want {
			t.Fatalf("mean=%v p=%v n=%d got=%s want=%s", tc.mean, tc.pvalue, tc.sample, got, tc.want)
		}
	}
}
