package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// from https://github.com/Anas099X/OpenSAT
// https://api.jsonsilo.com/public/942c3c3b-3a0c-4be3-81c2-12029def19f5
type OpenSATQuestion struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Visuals struct {
		Type       string `json:"type"`
		SvgContent string `json:"svg_content"`
	} `json:"visuals"`
	Question struct {
		Choices struct {
			A string `json:"A"`
			B string `json:"B"`
			C string `json:"C"`
			D string `json:"D"`
		} `json:"choices"`
		Question      string `json:"question"`
		Paragraph     string `json:"paragraph"`
		Explanation   string `json:"explanation"`
		CorrectAnswer string `json:"correct_answer"`
	} `json:"question"`
	Difficulty string `json:"difficulty"`
}

type Target struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Visuals struct {
		Type       string `json:"type"`
		SvgContent string `json:"svg_content"`
	} `json:"visuals"`
	Question struct {
		Choices       []string `json:"choices"`
		Question      string   `json:"question"`
		Paragraph     string   `json:"paragraph"`
		Explanation   string   `json:"explanation"`
		CorrectAnswer int      `json:"correct_answer"`
	} `json:"question"`
	Difficulty string `json:"difficulty"`
	Topic      string `json:"topic"`
}

func loadParsed(source map[string][]OpenSATQuestion) []Target {
	// Convert each question to Target format
	var targets []Target
	for topic, questions := range source {
		for _, question := range questions {
			target := convertToTarget(question)
			target.Topic = topic
			targets = append(targets, target)
		}
	}
	return targets
}

func loadOpenSAT(path string) (map[string][]OpenSATQuestion, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var source map[string][]OpenSATQuestion
	err = json.Unmarshal(file, &source)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}
	return source, nil
}

// convertToTarget converts an OpenSATQuestion to the Target format,
// shuffling the choices randomly on each call.
func convertToTarget(src OpenSATQuestion) Target {
	var target Target
	target.ID = src.ID
	target.Domain = src.Domain
	target.Visuals.Type = src.Visuals.Type
	target.Visuals.SvgContent = src.Visuals.SvgContent
	target.Difficulty = src.Difficulty

	// Copy question text, paragraph, and explanation
	target.Question.Question = src.Question.Question
	target.Question.Paragraph = src.Question.Paragraph
	target.Question.Explanation = src.Question.Explanation

	// Store original choices and the text of the correct answer
	originalChoices := []string{
		src.Question.Choices.A,
		src.Question.Choices.B,
		src.Question.Choices.C,
		src.Question.Choices.D,
	}
	correctAnswerIndex := letterToIndex(src.Question.CorrectAnswer)
	if correctAnswerIndex < 0 || correctAnswerIndex >= len(originalChoices) {
		// Handle invalid correct answer letter gracefully (e.g., log error, default to 0)
		fmt.Printf("Warning: Invalid correct answer '%s' for question ID %s. Defaulting to index 0.\n", src.Question.CorrectAnswer, src.ID)
		correctAnswerIndex = 0 // Or handle as appropriate
	}

	shuffled, newAnswerIndex := shuffleChoices(originalChoices, correctAnswerIndex)
	target.Question.Choices = shuffled
	target.Question.CorrectAnswer = newAnswerIndex

	return target
}

func shuffleChoices(choices []string, correctAnswerIndex int) ([]string, int) {
	// Keep track of the correct answer's index during the shuffle
	currentCorrectIndex := correctAnswerIndex

	// Shuffle the originalChoices slice in place, updating the correct index during swaps
	rand.Shuffle(len(choices), func(i, j int) {
		// Perform the swap
		choices[i], choices[j] = choices[j], choices[i]
		// Update the tracked index if the correct answer was involved in the swap
		if i == currentCorrectIndex {
			currentCorrectIndex = j
		} else if j == currentCorrectIndex {
			currentCorrectIndex = i
		}
	})
	return choices, currentCorrectIndex
}

func letterToIndex(letter string) int {
	switch strings.ToUpper(letter) {
	case "A":
		return 0
	case "B":
		return 1
	case "C":
		return 2
	case "D":
		return 3
	default:
		return -1 // indicates invalid answer
	}
}
