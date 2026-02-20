package polyaddr

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func HexToBytes(address string) ([]byte, error) {
	clean := strings.TrimPrefix(strings.ToLower(address), "0x")
	if len(clean) != 40 {
		return nil, fmt.Errorf("invalid address length")
	}
	buf, err := hex.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("decode address: %w", err)
	}
	return buf, nil
}

func BytesToHex(address []byte) string {
	return "0x" + hex.EncodeToString(address)
}
