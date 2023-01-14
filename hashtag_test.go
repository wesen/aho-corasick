package ahocorasick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func buildComplexTrie() *Trie {
	cleanerStrings := []string{
		"slon",
		"a",
		"scar",
		"et",
		"clean",
		"long",
		"scarp",
		"carpe",
		"cleaner",
		"leaner",
		"this",
		"s",
		"ane",
		"er",
		"i",
		"is",
		"carp",
		"scarpe",
		"an",
		"n",
		"le",
		"o",
		"cle",
		"th",
		"ar",
		"ean",
		"on",
	}
	return buildTrie(cleanerStrings)
}

func buildTrie(cleanerStrings []string) *Trie {
	builder := NewTrieBuilder()
	builder.AddStrings(cleanerStrings)
	return builder.Build()
}

func TestSingleWordMatches(t *testing.T) {
	trie := buildComplexTrie()
	s := "cleaner"
	trieMatches := trie.MatchString(s)
	matches := NewStringMatches(s, trieMatches)

	allMatches := matches.AllMatches
	assert.Equal(t, "cleaner", allMatches[0][0].MatchString())
	assert.Equal(t, "clean", allMatches[0][1].MatchString())
	assert.Equal(t, "leaner", allMatches[1][0].MatchString())

	for _, match := range matches.AllMatches {
		for _, m := range match {
			t.Logf("match: pos: %d - %s", m.Pos(), m.MatchString())
		}
	}

}

func TestSingleWordHashtags(t *testing.T) {
	trie := buildTrie([]string{"cleaner", "clean", "leaner"})

	s := "cleaner"
	trieMatches := trie.MatchString(s)
	matches := NewStringMatches(s, trieMatches)

	hashtags := matches.ComputeHashTags(0)
	// we expect something like:
	// - clean
	// - cLeaner
	// - cleanER
	for _, ht := range hashtags {
		t.Logf("hashtag: %s - %d", ht.String, ht.Score())
	}
}
