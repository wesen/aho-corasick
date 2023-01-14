package ahocorasick

import (
	"sort"
	"strings"
)

type StringMatches struct {
	String             string
	AllMatches         [][]*Match
	unprocessedMatches [][]*Match
	done               []bool
	cache              [][]*HashTag
}

func NewStringMatches(s string, matches []*Match) *StringMatches {
	matches_ := make([][]*Match, len(s))
	unProcessedMatches_ := make([][]*Match, len(s))
	done_ := make([]bool, len(s))

	for _, match := range matches {
		pos := match.Pos()

		if matches_[pos] == nil {
			matches_[pos] = make([]*Match, 0)
		}
		matches_[pos] = append(matches_[pos], match)
	}

	// we sort the individual matches to have the longest one first (most salient)
	for _, ms_ := range matches_ {
		sort.Slice(ms_, func(i, j int) bool {
			return len(ms_[i].Match()) > len(ms_[j].Match())
		})
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
