package main

import (
	"context"
	"log/slog"
	"math/rand"
	"sort"
	"time"
)

type DifficultyLevel string

const (
	Easy   DifficultyLevel = "Easy"
	Medium DifficultyLevel = "Medium"
	Hard   DifficultyLevel = "Hard"
)

type Difficulty struct {
	Level      DifficultyLevel
	Percentage float64
}

type TopicRequest struct {
	Count           int
	DifficultyLevel DifficultyLevel
}

func dailyRand(timestamp time.Time) *rand.Rand {
	// Get current date (year, month, day) as seed
	seed := timestamp.Year()*10000 + int(timestamp.Month())*100 + timestamp.Day()
	return rand.New(rand.NewSource(int64(seed)))
}

func generateRequests(topics map[string]int, day time.Weekday) map[string][]TopicRequest {
	difficulties := difficulty(day)
	if len(difficulties) < 1 {
		return nil
	}

	requests := make(map[string][]TopicRequest)

	for topic, total := range topics {
		remaining := total
		var distributed []TopicRequest

		// Distribute all but last difficulty
		// to prevent round up and round down errors
		for i := range len(difficulties) - 1 {
			d := difficulties[i]
			count := int(float64(total) * d.Percentage)
			distributed = append(distributed, TopicRequest{
				Count:           count,
				DifficultyLevel: d.Level,
			})
			remaining -= count
		}

		// Assign remaining to last difficulty
		last := difficulties[len(difficulties)-1]
		distributed = append(distributed, TopicRequest{
			Count:           remaining,
			DifficultyLevel: last.Level,
		})

		requests[topic] = distributed
	}

	return requests
}

func difficulty(day time.Weekday) []Difficulty {
	switch day {
	case time.Monday:
		return []Difficulty{{Level: Easy, Percentage: 1}}
	case time.Tuesday:
		return []Difficulty{
			{Level: Easy, Percentage: 0.5},   // 70% Easy
			{Level: Medium, Percentage: 0.5}, // 30% Medium
		}
	case time.Wednesday:
		return []Difficulty{
			{Level: Medium, Percentage: 1.0}, // 100% Medium
		}
	case time.Thursday:
		return []Difficulty{
			{Level: Medium, Percentage: 0.7}, // 70% Medium
			{Level: Hard, Percentage: 0.3},   // 30% Hard
		}
	case time.Friday:
		return []Difficulty{
			{Level: Medium, Percentage: 0.5}, // 50% Medium
			{Level: Hard, Percentage: 0.5},   // 50% Hard
		}
	case time.Saturday:
		return []Difficulty{
			{Level: Hard, Percentage: 1.0}, // 100% Hard
		}
	default:
		return []Difficulty{
			{Level: Easy, Percentage: 0.33},
			{Level: Medium, Percentage: 0.34},
			{Level: Hard, Percentage: 0.33},
		}
	}
}

// Pre-group questions by difficulty for topic
func groupByDifficulty(questions []OpenSATQuestion, requests []TopicRequest) map[string][]OpenSATQuestion {
	difficulties := make(map[string][]OpenSATQuestion)
	for _, question := range questions {
		// Only collect questions with difficulties we need
		needDifficulty := false
		for _, req := range requests {
			if string(req.DifficultyLevel) == question.Difficulty {
				needDifficulty = true
				break
			}
		}

		if needDifficulty {
			difficulties[question.Difficulty] = append(difficulties[question.Difficulty], question)
		}
	}
	return difficulties
}

func shuffleSubset(ctx context.Context, allQuestions map[string][]OpenSATQuestion, topicCounts map[string]int) map[string][]Target {
	now := time.Now()
	rnd := dailyRand(now)
	topics := make(map[string][]Target, len(topicCounts))
	requests := generateRequests(topicCounts, now.Weekday())

	// Get topic keys and sort them for deterministic order
	topicKeys := make([]string, 0, len(allQuestions))
	for k := range allQuestions {
		topicKeys = append(topicKeys, k)
	}
	sort.Strings(topicKeys)

	// Iterate over sorted keys
	for _, topic := range topicKeys {
		if _, exists := requests[topic]; !exists {
			// Skip topics that aren't in our requests
			continue
		}

		difficulties := groupByDifficulty(allQuestions[topic], requests[topic])

		// Preallocate targets slice based on total count needed
		totalCount := 0
		for _, req := range requests[topic] {
			totalCount += req.Count
		}
		targetQuestions := make([]Target, 0, totalCount)

		// Process each difficulty request
		for _, request := range requests[topic] {
			diffQuestions := difficulties[string(request.DifficultyLevel)]
			count := min(request.Count, len(diffQuestions))
			if count < 0 {
				break
			}

			if count > len(diffQuestions)/2 {
				slog.WarnContext(ctx, "Requesting large number of questions. Should use Fisher-Yates", "count", count)
			}

			// Select random subset without full shuffle
			// Use a map to track selected indices for small subsets
			selected := make(map[int]bool, count)
			for range count {
				// Find an unselected index
				var idx int
				for {
					idx = rnd.Intn(len(diffQuestions))
					if !selected[idx] {
						selected[idx] = true
						break
					}
				}
				targetQuestions = append(targetQuestions, convertToTarget(diffQuestions[idx]))
			}
		}

		topics[topic] = targetQuestions
	}

	return topics
}
