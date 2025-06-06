package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
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
	Difficulty DifficultyLevel `json:"difficulty"`
	Topic      string          `json:"topic"`
}

//go:embed OpenSAT.json
var openSAT embed.FS

func loadOpenSAT() (map[string][]OpenSATQuestion, error) {
	file, err := openSAT.ReadFile("OpenSAT.json")
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

func listStats(stats map[string]int) {
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("\t\t%s: %d\n", key, stats[key])
	}
}

func showStats(allQuestions map[string][]OpenSATQuestion) {
	for topic, questions := range allQuestions {
		fmt.Printf("%s: %d\n", topic, len(questions))
		domains := map[string]int{}
		difficulties := map[string]int{}
		breakdown := map[string]int{}
		for _, question := range questions {
			if question.Visuals.SvgContent != "" {
				fmt.Println(question.ID)
			}
			domains[question.Domain]++
			difficulties[question.Difficulty]++
			breakdown[fmt.Sprintf("%s_%s", question.Domain, question.Difficulty)]++
		}

		fmt.Println("\tDOMAINS")
		listStats(domains)
		fmt.Println("\tDIFFICULTIES")
		listStats(difficulties)
		fmt.Println("\tBREAKDOWN")
		listStats(breakdown)
	}
}

// convertToTarget converts an OpenSATQuestion to the Target format,
func convertToTarget(src OpenSATQuestion, topic string) Target {
	var target Target
	target.ID = src.ID
	target.Domain = src.Domain
	target.Visuals.Type = src.Visuals.Type
	target.Visuals.SvgContent = src.Visuals.SvgContent
	target.Difficulty = DifficultyLevel(src.Difficulty)
	target.Question.Question = src.Question.Question
	target.Question.Paragraph = src.Question.Paragraph
	target.Question.Explanation = src.Question.Explanation
	target.Topic = topic

	target.Question.Choices = []string{
		src.Question.Choices.A,
		src.Question.Choices.B,
		src.Question.Choices.C,
		src.Question.Choices.D,
	}
	target.Question.CorrectAnswer = letterToIndex(src.Question.CorrectAnswer)
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
	if len(letter) == 0 {
		return -1 // Empty string
	}

	// Get the first character and convert to uppercase if needed
	char := letter[0]
	if char > 'D' {
		// invalid letter
		return -1
	}
	if char >= 'a' && char <= 'd' {
		char -= 32 // Convert lowercase to uppercase (ASCII difference)
	}

	// Direct mapping from character to index
	if char >= 'A' {
		return int(char - 'A')
	}

	return -1 // Invalid answer
}
