package unicode

import (
	"github.com/rivo/uniseg"
)

// CharWidth calculates the display width of a single character
// This handles emoji, CJK characters, and other wide characters correctly
func CharWidth(r rune) int {
	return uniseg.StringWidth(string(r))
}

// StringDisplayWidth calculates the total display width of a string
// This properly handles grapheme clusters and wide characters
func StringDisplayWidth(s string) int {
	return uniseg.StringWidth(s)
}

// GraphemeCount returns the number of user-perceived characters (graphemes) in a string
// This is different from rune count for complex Unicode like emoji with modifiers
func GraphemeCount(s string) int {
	g := uniseg.NewGraphemes(s)
	count := 0
	for g.Next() {
		count++
	}
	return count
}

// TruncateToWidth truncates a string to fit within a specified display width
// Returns the truncated string and the actual width used
func TruncateToWidth(s string, maxWidth int) (string, int) {
	if maxWidth <= 0 {
		return "", 0
	}

	g := uniseg.NewGraphemes(s)
	var result string
	currentWidth := 0
	
	for g.Next() {
		grapheme := g.Str()
		graphemeWidth := uniseg.StringWidth(grapheme)
		
		if currentWidth+graphemeWidth > maxWidth {
			break
		}
		
		result += grapheme
		currentWidth += graphemeWidth
	}
	
	return result, currentWidth
}

// PadToWidth pads a string with spaces to reach the specified display width
func PadToWidth(s string, targetWidth int) string {
	currentWidth := StringDisplayWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	
	padding := targetWidth - currentWidth
	for i := 0; i < padding; i++ {
		s += " "
	}
	return s
}

// IsWideChar determines if a character is considered "wide" (takes 2 columns)
func IsWideChar(r rune) bool {
	return CharWidth(r) == 2
}