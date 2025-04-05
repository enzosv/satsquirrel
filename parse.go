package main

import (
	"encoding/json"
	"fmt"
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

	// Convert choices from map to ordered slice
	target.Question.Choices = []string{
		src.Question.Choices.A,
		src.Question.Choices.B,
		src.Question.Choices.C,
		src.Question.Choices.D,
	}

	// Convert correct answer from letter to index
	target.Question.CorrectAnswer = letterToIndex(src.Question.CorrectAnswer)

	return target
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
