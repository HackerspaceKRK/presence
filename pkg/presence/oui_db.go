package presence

import (
	"bytes"
	"encoding/hex"
	"strings"
)

//go:generate bash -c "go run ./data/gen_oui.go ./data/oui.txt > oui_db_generated.go"

type ouiDbKeyType string

var OuiDbKey ouiDbKeyType = "ouiDb"

type OuiRecord struct {
	Prefix [3]byte
	Vendor string
}

// LookupOuiByPrefix performs a binary search on ouiRecords to find a record by prefix.
func LookupOuiByPrefix(prefix [3]byte) (OuiRecord, bool) {
	low, high := 0, len(ouiRecords)-1
	for low <= high {
		mid := (low + high) / 2
		cmp := bytes.Compare(ouiRecords[mid].Prefix[:], prefix[:])
		if cmp == 0 {
			return ouiRecords[mid], true
		} else if cmp < 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return OuiRecord{}, false
}

// LookupOuiByMAC returns the vendor name based on a MAC address or an empty string if not found.
func LookupOuiByMAC(macAddress string) string {
	normalizedMAC := strings.ReplaceAll(macAddress, ":", "")
	normalizedMAC = strings.ReplaceAll(normalizedMAC, "-", "")
	normalizedMAC = strings.ToUpper(normalizedMAC)

	// Extract the prefix from the MAC address (first 6 characters)
	prefixHex := normalizedMAC[:6]

	prefixBytes, err := hex.DecodeString(prefixHex)
	if err != nil {
		return ""
	}

	var prefix [3]byte
	copy(prefix[:], prefixBytes)

	if record, found := LookupOuiByPrefix(prefix); found {
		return record.Vendor
	}
	return ""
}
