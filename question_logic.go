package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Difficulty struct {
	Level      int
	Percentage float64
}

type TopicRequest struct {
	Topic      string
	Count      int
	Difficulty string
}

func dailyRand(timestamp time.Time) *rand.Rand {
	// Get current date (year, month, day) as seed
	seed := timestamp.Year()*10000 + int(timestamp.Month())*100 + timestamp.Day()
	fmt.Println("seed", seed)
	return rand.New(rand.NewSource(int64(seed)))
}

func difficulty(day time.Weekday) []Difficulty {
	switch day {
	case time.Monday:
		return []Difficulty{{Level: 1, Percentage: 1}}
	case time.Tuesday:
		return []Difficulty{
			{Level: 1, Percentage: 0.5}, // 70% Easy
			{Level: 2, Percentage: 0.5}, // 30% Medium
		}
	case time.Wednesday:
		return []Difficulty{
			{Level: 2, Percentage: 1.0}, // 100% Medium
		}
	case time.Thursday:
		return []Difficulty{
			{Level: 2, Percentage: 0.7}, // 70% Medium
			{Level: 3, Percentage: 0.3}, // 30% Hard
		}
	case time.Friday:
		return []Difficulty{
			{Level: 2, Percentage: 0.5}, // 50% Medium
			{Level: 3, Percentage: 0.5}, // 50% Hard
		}
	case time.Saturday:
		return []Difficulty{
			{Level: 3, Percentage: 1.0}, // 100% Hard
		}
	default:
		return []Difficulty{
			{Level: 1, Percentage: 0.33},
			{Level: 2, Percentage: 0.34},
			{Level: 3, Percentage: 0.33},
		}
	}
}

func randomize(allQuestions map[string][]OpenSATQuestion, topicCounts map[string]int) map[string][]Target {
	now := time.Now()

	rnd := dailyRand(now)
	topics := map[string][]Target{}
	n := 0 // To avoid shadowing in the loop below

	for topic, questions := range allQuestions {

		difficulties := map[string][]OpenSATQuestion{}
		for _, question := range questions {
			if question.Difficulty != "Easy" {
				continue
			}
			difficulties[question.Difficulty] = append(difficulties[question.Difficulty], question)
		}

		n = len(difficulties["Easy"])
		count := topicCounts[topic]
		if count > n {
			fmt.Printf("Warning: Requested %d questions for topic '%s', but only %d available. Returning all available.\n", count, topic, n)
			count = n
		}

		// Allocate target slice directly
		targetQuestions := make([]Target, count)

		// Perform partial Fisher-Yates shuffle, converting and assigning directly
		for i := range count {
			// Choose index j from the remaining part [i, n-1]
			j := i + rnd.Intn(n-i)
			fmt.Println(topic, j)
			// Swap elements in the original slice
			difficulties["Easy"][i], difficulties["Easy"][j] = difficulties["Easy"][j], difficulties["Easy"][i]
			// Convert the element now at index i (which came from index j)
			// and place it directly into the target slice
			targetQuestions[i] = convertToTarget(difficulties["Easy"][i])

			// not targetQuestions[i] = convertToTarget(difficulties["easy"][j]) to avoid duplicates
		}
		fmt.Println()
		topics[topic] = targetQuestions
	}

	return topics
}
