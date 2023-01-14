package main

import (
	"github.com/BobuSumisu/aho-corasick/cmd/hashtag/cmds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashtag",
	Short: "hashtag is a tool for finding hashtags in text",
}

func init() {
	rootCmd.AddCommand(cmds.ReplCmd)
	wordLists := []string{
		"test_data/words",
		"test_data/words.txt",
		"test_data/google-10000-english-no-swears.txt",
	}
	rootCmd.PersistentFlags().StringSlice("dict", wordLists, "Dictionary file(s) to use")
}

func main() {
	_ = rootCmd.Execute()

}
