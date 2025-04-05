package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Difficulty struct {
	Level      int
	Percentage float64
}

func (d Difficulty) difficultyString() string {
	switch d.Level {
	case 1:
		return "Easy"
	case 2:
		return "Medium"
	case 3:
		return "Hard"
	default:
		return "Unknown"
	}
}

type TopicRequest struct {
	Count      int
	Difficulty string
}

func dailyRand(timestamp time.Time) *rand.Rand {
	// Get current date (year, month, day) as seed
	seed := timestamp.Year()*10000 + int(timestamp.Month())*100 + timestamp.Day()
	fmt.Println("seed", seed)
	return rand.New(rand.NewSource(int64(seed)))
}

func generateRequests(topics map[string]int, day time.Weekday) map[string][]TopicRequest {
	requests := map[string][]TopicRequest{}
	difficulties := difficulty(day)
	for topic, count := range topics {
		for _, difficulty := range difficulties {
			requests[topic] = append(requests[topic], TopicRequest{int(difficulty.Percentage) * count, difficulty.difficultyString()})
		}
	}
	return requests
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

	requests := generateRequests(topicCounts, now.Weekday())

	// Get topic keys and sort them for deterministic order
	topicKeys := make([]string, 0, len(allQuestions))
	for k := range allQuestions {
		topicKeys = append(topicKeys, k)
	}
	sort.Strings(topicKeys)

	// Iterate over sorted keys
	for _, topic := range topicKeys {
		questions := allQuestions[topic] // Get questions for the current topic

		difficulties := map[string][]OpenSATQuestion{}
		for _, question := range questions {
			// TODO: if difficulty not in request[topic], skip
			difficulties[question.Difficulty] = append(difficulties[question.Difficulty], question)
		}

		// Allocate target slice directly
		var targetQuestions []Target

		for _, request := range requests[topic] {

			// Perform partial Fisher-Yates shuffle, converting and assigning directly
			for i := range request.Count {
				n := len(difficulties[request.Difficulty])
				// Choose index j from the remaining part [i, n-1]
				j := i + rnd.Intn(n-i)
				fmt.Println(topic, j)
				// Swap elements in the original slice
				difficulties[request.Difficulty][i], difficulties[request.Difficulty][j] = difficulties[request.Difficulty][j], difficulties[request.Difficulty][i]
				// Convert the element now at index i (which came from index j)
				// and place it directly into the target slice
				targetQuestions = append(targetQuestions, convertToTarget(difficulties[request.Difficulty][i]))

				// not targetQuestions[i] = convertToTarget(difficulties["easy"][j]) to avoid duplicates
			}
			fmt.Println()
		}
		topics[topic] = targetQuestions
	}

	return topics
}
