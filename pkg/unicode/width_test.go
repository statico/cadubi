package unicode

import (
	"testing"
)

func TestCharWidth(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected int
	}{
		{"ASCII letter", 'A', 1},
		{"Space", ' ', 1},
		{"Number", '1', 1},
		{"Emoji smiley", '😀', 2},
		{"Heart emoji", '❤', 1}, // Actually width 1 in terminals
		{"CJK character", '中', 2},
		{"Box drawing", '┌', 1},
		{"Arrow", '→', 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CharWidth(tt.char)
			if result != tt.expected {
				t.Errorf("CharWidth(%c) = %d; want %d", tt.char, result, tt.expected)
			}
		})
	}
}

func TestStringDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII only", "Hello", 5},
		{"With emoji", "Hi 😀", 5}, // H(1) + i(1) + space(1) + emoji(2) = 5
		{"With emoji corrected", "Hi😀", 4}, // H(1) + i(1) + emoji(2) = 4
		{"CJK text", "中文", 4},
		{"Mixed", "A中B", 4}, // A(1) + 中(2) + B(1) = 4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringDisplayWidth(tt.input)
			if result != tt.expected {
				t.Errorf("StringDisplayWidth(%q) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsWideChar(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		{"ASCII letter", 'A', false},
		{"Emoji smiley", '😀', true},
		{"CJK character", '中', true},
		{"Box drawing", '┌', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWideChar(tt.char)
			if result != tt.expected {
				t.Errorf("IsWideChar(%c) = %t; want %t", tt.char, result, tt.expected)
			}
		})
	}
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		maxWidth      int
		expectedStr   string
		expectedWidth int
	}{
		{"Empty string", "", 5, "", 0},
		{"Short ASCII", "Hi", 5, "Hi", 2},
		{"Exact width", "Hello", 5, "Hello", 5},
		{"Too long ASCII", "Hello World", 5, "Hello", 5},
		{"With emoji", "Hi😀", 3, "Hi", 2}, // Can't fit emoji
		{"With emoji fits", "Hi😀", 4, "Hi😀", 4},
		{"CJK truncate", "中文字", 3, "中", 2}, // Only first char fits
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultStr, resultWidth := TruncateToWidth(tt.input, tt.maxWidth)
			if resultStr != tt.expectedStr || resultWidth != tt.expectedWidth {
				t.Errorf("TruncateToWidth(%q, %d) = (%q, %d); want (%q, %d)", 
					tt.input, tt.maxWidth, resultStr, resultWidth, tt.expectedStr, tt.expectedWidth)
			}
		})
	}
}