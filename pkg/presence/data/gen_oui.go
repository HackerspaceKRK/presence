package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type OuiRecord struct {
	Prefix [3]byte
	Vendor string
}

func main() {
	// Open the file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var records []OuiRecord
	scanner := bufio.NewScanner(file)
	hexRegex := regexp.MustCompile(`^[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}`)
	for scanner.Scan() {
		line := scanner.Text()
		matched := hexRegex.MatchString(line)

		if matched {
			ouiHex := strings.Split(line, " ")[0]

			ouiHex = strings.ReplaceAll(ouiHex, "-", "")
			var prefix [3]byte
			decoded, err := hex.DecodeString(ouiHex)
			if err != nil {
				panic(err)
			}
			copy(prefix[:], decoded)

			vendor := strings.TrimSpace(line)
			vendor = strings.ReplaceAll(vendor, "\t", " ")
			vendor = strings.ReplaceAll(vendor, "(base 16)", "")
			vendor = strings.ReplaceAll(vendor, "(hex)", "")
			vendor = hexRegex.ReplaceAllString(vendor, "")
			vendor = strings.TrimSpace(vendor)
			records = append(records, OuiRecord{
				Prefix: prefix,
				Vendor: vendor,
			})

		}
	}

	// Sort records by prefix
	sort.Slice(records, func(i, j int) bool {
		return bytes.Compare(records[i].Prefix[:], records[j].Prefix[:]) < 0
	})

	// Generate the Go code
	code := "package presence\n\n"
	code += "var ouiRecords = []OuiRecord{\n"
	for _, record := range records {
		code += fmt.Sprintf("\t{[3]byte{0x%02x, 0x%02x, 0x%02x}, \"%s\"},\n", record.Prefix[0], record.Prefix[1], record.Prefix[2], record.Vendor)
	}
	code += "}\n"

	// Output the generated Go code
	fmt.Println(code)
}
