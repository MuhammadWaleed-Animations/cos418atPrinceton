package cos418_hw1_1

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
	"fmt"
)

// Find the top K most common words in a text document.
func topWords(path string, numWords int, charThreshold int) []WordCount {
	// Open the file
	file, err := os.Open(path)
	checkError(err)
	defer file.Close()

	// Regular expression to keep only alphanumeric characters
	re := regexp.MustCompile(`[^0-9a-zA-Z]+`)

	wordFreq := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line) // Split line into words
		for _, word := range words {
			// Convert to lowercase and remove non-alphanumeric characters
			cleanWord := re.ReplaceAllString(strings.ToLower(word), "")
			if len(cleanWord) >= charThreshold {
				wordFreq[cleanWord]++
			}
		}
	}

	checkError(scanner.Err())

	// Convert map to a slice of WordCount
	wordCounts := make([]WordCount, 0, len(wordFreq))
	for word, count := range wordFreq {
		wordCounts = append(wordCounts, WordCount{Word: word, Count: count})
	}

	// Sort the words
	sortWordCounts(wordCounts)

	// Return the top numWords words
	if numWords > len(wordCounts) {
		numWords = len(wordCounts)
	}
	return wordCounts[:numWords]
}


// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.
// DO NOT MODIFY THIS FUNCTION!
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
