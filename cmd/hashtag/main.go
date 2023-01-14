package main

import (
	"fmt"
	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"runtime"
	"time"
)

func main() {
	builder := ahocorasick.NewTrieBuilder()
	fmt.Printf("Loading dictionary...\n")
	wordLists := []string{
		//"test_data/words",
		"test_data/words.txt",
		//"test_data/google-10000-english-no-swears.txt",
	}
	var err error
	for _, wordList := range wordLists {
		fmt.Printf("Loading %s...\n", wordList)
		err = builder.LoadStrings(wordList)

		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Building trie...\n")
	trie := builder.Build()
	fmt.Printf("Built.\n")

	// pirint allocated memory size from garbage collector information
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Allocated memory size: %s\n", formatMemorySize(mem.Alloc))

	// read strings from stdin
	// for each string, find all matches
	for {
		var s string
		_, err := fmt.Scanln(&s)
		if err != nil {
			break
		}
		//s = "this"
		//s = "thisisatest"

		// we keep a list of all matches starting at a certain position

		start := time.Now()
		var hashTags []*ahocorasick.HashTag
		var trieMatches []*ahocorasick.Match
		iterCount := 1
		for i := 0; i < iterCount; i++ {
			trieMatches = trie.MatchString(s)
		}
		elapsed := time.Since(start)
		fmt.Printf("Aho corasick took %d ns for %d matches\n",
			elapsed.Nanoseconds()/int64(iterCount),
			len(trieMatches),
		)

		matchedStrings := make(map[string]interface{})
		for _, m := range trieMatches {
			fmt.Printf(" pos: %d - %s\n", m.Pos(), m.String())
			matchedStrings[string(m.Match())] = nil
		}
		for k := range matchedStrings {
			fmt.Printf("\"%s\",\n", k)
		}

		start = time.Now()
		iterCount = 1
		for i := 0; i < iterCount; i++ {
			matches := ahocorasick.NewStringMatches(s, trieMatches)
			hashTags = matches.SuggestHashtags()
		}
		elapsed = time.Since(start)
		fmt.Printf("Hashtag took %d ns for %d hashtags\n",
			elapsed.Nanoseconds()/int64(iterCount),
			len(hashTags),
		)

		// show at most 5 results
		for _, hashTag := range hashTags[:5] {
			fmt.Printf("%d - %s\n", hashTag.Words, hashTag.String)
		}
	}
}

func formatMemorySize(alloc uint64) string {
	// convert to GB
	allocGB := float64(alloc) / 1024 / 1024 / 1024
	return fmt.Sprintf("%.2f GB", allocGB)
}
