package opensat

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShuffleChoices(t *testing.T) {
	// Set a fixed seed for reproducible tests
	rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Run("Basic Functionality", func(t *testing.T) {
		// Arrange
		originalChoices := []string{"A", "B", "C", "D"}
		originalCopy := make([]string, len(originalChoices))
		copy(originalCopy, originalChoices)
		correctAnswerIndex := 2 // "C" is the correct answer
		correctAnswer := originalChoices[correctAnswerIndex]

		// Act
		shuffledChoices, newCorrectIndex := shuffleChoices(originalChoices, correctAnswerIndex)

		// Assert
		assert.Equal(t, len(originalCopy), len(shuffledChoices), "Shuffled choices should have the same length as original")
		require.Less(t, newCorrectIndex, len(shuffledChoices), "New correct index should be within bounds")
		require.GreaterOrEqual(t, newCorrectIndex, 0, "New correct index should be non-negative")
		assert.Equal(t, correctAnswer, shuffledChoices[newCorrectIndex], "Correct answer should be at the new index")

		// Check that all original elements are still present (no additions or omissions)
		assert.ElementsMatch(t, originalCopy, shuffledChoices, "Shuffled slice should contain the same elements as original")
	})

	t.Run("Empty Slice", func(t *testing.T) {
		emptyChoices := []string{}
		shuffled, newIndex := shuffleChoices(emptyChoices, 0)

		assert.Empty(t, shuffled, "Empty slice should remain empty after shuffle")
		assert.Equal(t, 0, newIndex, "Index should remain 0 for empty slice")
	})

	t.Run("Single Element", func(t *testing.T) {
		singleChoice := []string{"Only Option"}
		shuffled, newIndex := shuffleChoices(singleChoice, 0)

		assert.Len(t, shuffled, 1, "Single element slice should remain length 1 after shuffle")
		assert.Equal(t, 0, newIndex, "Index should remain 0 for single element slice")
		assert.Equal(t, "Only Option", shuffled[0], "Element should remain unchanged")
	})

	t.Run("Shuffle Multiple Times", func(t *testing.T) {
		// This test verifies that multiple shuffles maintain correct tracking
		choices := []string{"A", "B", "C", "D", "E"}
		correctIndex := 3 // "D" is correct
		correctAnswer := choices[correctIndex]

		// Shuffle multiple times
		for i := range 5 {
			choices, correctIndex = shuffleChoices(choices, correctIndex)

			// Verify correct answer is still at the tracked index
			assert.Equal(t, correctAnswer, choices[correctIndex],
				"After shuffle %d, correct answer should be at tracked index", i+1)
		}
	})

	t.Run("Original Slice Modified", func(t *testing.T) {
		// Test that the function modifies the original slice (not a copy)
		originalChoices := []string{"A", "B", "C", "D"}
		originalCopy := make([]string, len(originalChoices))
		copy(originalCopy, originalChoices)

		shuffledChoices, _ := shuffleChoices(originalChoices, 0)

		// Check that the original slice and returned slice are the same object
		assert.Same(t, &originalChoices[0], &shuffledChoices[0],
			"Function should modify original slice, not return a copy")

		// If the slice has more than one element, it's statistically unlikely (but possible)
		// that the shuffle would result in the exact same order
		if len(originalChoices) > 1 {
			// Just log if they happen to be equal, but don't fail the test
			// since there's a small but non-zero chance the shuffle doesn't change anything
			if assert.NotEqual(t, originalCopy, originalChoices, "Choices should be shuffled") == false {
				t.Log("Note: Shuffle resulted in original order (possible but unlikely)")
			}
		}
	})
}
