package utils

import (
	"strings"
	"testing"
)

// TestGenerateShortID tests the GenerateShortID function to ensure it produces unique IDs of the correct format.
func TestGenerateShortID(t *testing.T) {
	idMap := make(map[string]bool)
	count := 1000 // Number of IDs to generate

	for i := 0; i < count; i++ {
		id := GenerateShortID()
		// Check for correct length
		if len(id) != 8 {
			t.Errorf("ID '%s' does not have the expected length of 8 characters", id)
		}
		// Check for hyphens
		if strings.Contains(id, "-") {
			t.Errorf("ID '%s' contains a hyphen", id)
		}
		// Check for uniqueness
		if _, exists := idMap[id]; exists {
			t.Errorf("ID '%s' is not unique", id)
		}
		idMap[id] = true
	}
}
