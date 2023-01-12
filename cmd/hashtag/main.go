package main

import (
	"fmt"
	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"sort"
	"strings"
	"time"
)

type StringMatches struct {
	String             string
	AllMatches         [][]*ahocorasick.Match
	unprocessedMatches [][]*ahocorasick.Match
	done               []bool
	cache              [][]*HashTag
}

func NewStringMatches(s string, matches []*ahocorasick.Match) *StringMatches {
	matches_ := make([][]*ahocorasick.Match, len(s))
	unProcessedMatches_ := make([][]*ahocorasick.Match, len(s))
	done_ := make([]bool, len(s))

	for _, match := range matches {
		pos := match.Pos()

		if matches_[pos] == nil {
			matches_[pos] = make([]*ahocorasick.Match, 0)
		}
		matches_[pos] = append(matches_[pos], match)
	}

	for i, m := range matches_ {
		unProcessedMatches_[i] = m
	}

	return &StringMatches{
		s,
		matches_,
		unProcessedMatches_,
		done_,
		make([][]*HashTag, len(s)),
	}
}

type HashTag struct {
	String string
	Words  int
}

func (ht *HashTag) Score() int {
	return ht.Words
}

func (sm *StringMatches) ComputeHashTags(pos int) []*HashTag {
	//fmt.Printf("Computing hashtags for %d: %s\n", pos, sm.String[pos:])

	if pos >= len(sm.String) {
		return []*HashTag{
			{
				String: "",
				Words:  0,
			},
		}
	}

	if sm.done[pos] {
		//fmt.Println("Already done")
		return sm.cache[pos]
	}

	if sm.cache[pos] == nil {
		sm.cache[pos] = make([]*HashTag, 0)
	}

	// we first start by looking at potential words at this position
	if sm.unprocessedMatches[pos] == nil || len(sm.unprocessedMatches[pos]) == 0 {
		// if we have no more matches to process, we can just gobble the character up
		// in case there is an unknown hashtag here, we just won't capitalize
		for _, ht := range sm.ComputeHashTags(pos + 1) {
			ht_ := &HashTag{
				sm.String[:pos] + capitalize(ht.String),
				// we increment the word despite not having a word here
				// because that will downweight the result
				ht.Words + 1,
			}
			//fmt.Printf("Adding %s to %s\n", ht_.String, sm.String[pos:pos+1])
			sm.cache[pos] = append(sm.cache[pos], ht_)
		}
		sm.done[pos] = true
		return sm.cache[pos]
	}

	// go over matches at the given position
	for _, match := range sm.unprocessedMatches[pos] {
		// we can now recurse
		s := match.Match()
		pos_ := match.Pos()
		matchLen := len(s)
		nextWordPos := int(pos_) + matchLen
		for _, ht := range sm.ComputeHashTags(nextWordPos) {
			ht_ := &HashTag{
				sm.String[pos:nextWordPos] + capitalize(ht.String),
				ht.Words + 1,
			}
			//fmt.Printf("Adding %s to %s\n", ht_.String, sm.String[pos:nextWordPos])
			sm.cache[pos] = append(sm.cache[pos], ht_)
		}
	}

	sm.done[pos] = true
	return sm.cache[pos]
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// SuggestHashtags using a DP approach to computing possible hashtags
// It keeps track of the best result starting at a certain position.
// A best hashtag is the one that uses the least capitalizations to cover a given area.
func (sm *StringMatches) SuggestHashtags() []*HashTag {
	hashTags := sm.ComputeHashTags(0)

	// sort hashTags by Words
	sort.Slice(hashTags, func(i, j int) bool {
		return hashTags[i].Words < hashTags[j].Words
	})

	return hashTags
}

func main() {
	builder := ahocorasick.NewTrieBuilder()
	fmt.Printf("Loading dictionary...\n")
	wordLists := []string{
		//"test_data/words",
		//"test_data/words.txt",
		"test_data/google-10000-english-no-swears.txt",
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
		var hashTags []*HashTag
		var trieMatches []*ahocorasick.Match
		iterCount := 100
		for i := 0; i < iterCount; i++ {
			trieMatches = trie.MatchString(s)
		}
		elapsed := time.Since(start)
		fmt.Printf("Aho corasick took %d ns for %d matches\n",
			elapsed.Nanoseconds()/int64(iterCount),
			len(trieMatches),
		)

		for _, m := range trieMatches {
			fmt.Printf(" pos: %d - %s\n", m.Pos(), m.String())
		}

		start = time.Now()
		iterCount = 1
		for i := 0; i < iterCount; i++ {
			matches := NewStringMatches(s, trieMatches)
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
