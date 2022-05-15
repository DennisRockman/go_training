package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topWordsQuantity, dash = 10, "-"

var reg = regexp.MustCompile(`[\p{L}\d-]+`)

func Top10(s string) []string {
	if len(s) == 0 {
		return nil
	}

	s = strings.ToLower(s)
	words := reg.FindAllString(s, -1)
	wordsWithCounters := make(map[string]int)
	for _, word := range words {
		wordsWithCounters[word]++
	}

	topWords := make([]string, 0, len(wordsWithCounters))
	wordsWithCounters[dash] = 0
	for k := range wordsWithCounters {
		topWords = append(topWords, k)
	}

	sort.Slice(topWords, func(i, j int) bool {
		return wordsWithCounters[topWords[i]] > wordsWithCounters[topWords[j]] ||
			(wordsWithCounters[topWords[i]] == wordsWithCounters[topWords[j]] && topWords[i] < topWords[j])
	})
	return topWords[:min(topWordsQuantity, len(topWords)-1)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
