package xrpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKeyName(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
		wantErr  bool
	}{
		{
			name:     "standard toml file",
			filePath: "/path/to/key.toml",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "file no extension",
			filePath: "/path/to/key",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "relative path",
			filePath: "key.toml",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "empty filename",
			filePath: "/path/to/.toml",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractKeyName(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetXRPLKeyDir(t *testing.T) {
	home := "/home/user"
	chain := "xrpl-mainnet"
	expected := []string{"/home/user", "keys", "xrpl-mainnet", "metadata"}

	got := getXRPLKeyDir(home, chain)
	assert.Equal(t, expected, got)
}

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
			got, err := StringToHex(tt.input, tt.length)
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
			gotBase, gotQuote, err := ParseAssetsFromSignal(tt.signalID)
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
			got, err := Uint64StrToHexStr(tt.input)
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
