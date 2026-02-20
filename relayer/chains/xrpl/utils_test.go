package xrpl_test

import (
	"testing"

	"github.com/bandprotocol/falcon/relayer/chains/xrpl"
)

// TestStringToHex tests hex encoding with optional length padding
func TestStringToHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		length  int
		want    string
		wantErr bool
	}{
		{"Standard conversion", "USD", 0, "555344", false},
		{"With padding", "USD", 40, "5553440000000000000000000000000000000000", false},
		{"Length too short", "LONGSTRING", 4, "", true},
		{"Empty string", "", 0, "", false},
		{"Empty string with padding", "", 4, "0000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xrpl.StringToHex(tt.input, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseAssetsFromSignal tests asset extraction from XRPL signal strings
func TestParseAssetsFromSignal(t *testing.T) {
	tests := []struct {
		name      string
		signalID  string
		wantBase  string
		wantQuote string
		wantErr   bool
	}{
		{
			name:      "Valid standard 3-char assets",
			signalID:  "CS:BTC-USD",
			wantBase:  "BTC",
			wantQuote: "USD",
			wantErr:   false,
		},
		{
			name:      "More than 3 chars asset",
			signalID:  "CS:RLUSD-USD",
			wantBase:  "524C555344000000000000000000000000000000",
			wantQuote: "USD",
			wantErr:   false,
		},
		{
			name:      "XRP special case",
			signalID:  "CS:XRP-USD",
			wantBase:  "XRP",
			wantQuote: "USD",
			wantErr:   false,
		},
		{
			name:     "Invalid format - missing dash",
			signalID: "CS:XRPUSD",
			wantErr:  true,
		},
		{
			name:     "Invalid format - empty asset",
			signalID: "CS:XRP- ",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBase, gotQuote, err := xrpl.ParseAssetsFromSignal(tt.signalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAssetsFromSignal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBase != tt.wantBase || gotQuote != tt.wantQuote {
				t.Errorf("ParseAssetsFromSignal() got = (%v, %v), want (%v, %v)", gotBase, gotQuote, tt.wantBase, tt.wantQuote)
			}
		})
	}
}

// TestUint64StrToHexStr tests conversion of large numeric strings to fixed-width hex
func TestUint64StrToHexStr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Small number", "255", "00000000000000FF", false},
		{"Large uint64 max", "18446744073709551615", "FFFFFFFFFFFFFFFF", false},
		{"Zero", "0", "0000000000000000", false},
		{"Invalid non-numeric", "abc", "", true},
		{"Negative", "-1", "", true},
		{"Large uint64 max + 1", "18446744073709551616", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xrpl.Uint64StrToHexStr(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64StrToHexStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64StrToHexStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
