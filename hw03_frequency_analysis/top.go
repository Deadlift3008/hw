package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const limit = 10

func Top10(text string) []string {
	words := strings.Fields(text)

	dict := make(map[string]int, len(words))

	for _, v := range words {
		currentCount := dict[v]
		dict[v] = currentCount + 1
	}

	result := make([]string, 0)

	for key := range dict {
		result = append(result, key)
	}

	sort.Slice(result, func(i, j int) bool {
		first := result[i]
		second := result[j]

		firstCount := dict[first]
		secondCount := dict[second]

		if firstCount == secondCount {
			return first < second
		}

		return firstCount > secondCount
	})

	if len(result) > limit {
		return result[:limit]
	}

	return result
}
