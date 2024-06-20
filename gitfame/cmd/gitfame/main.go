//go:build !solution

package main

import (
	"os"

	flag "github.com/spf13/pflag"

	"gitlab.com/slon/shad-go/gitfame/internal"
)

func main() {
	set := FlagSet()

	// Fetch command line arguments
	repository, _ := set.GetString("repository")
	revision, _ := set.GetString("revision")
	extensions, _ := set.GetStringSlice("extensions")
	languages, _ := set.GetStringSlice("languages")
	exclude, _ := set.GetStringSlice("exclude")
	restrictTo, _ := set.GetStringSlice("restrict-to")
	useCommitter, _ := set.GetBool("use-committer")
	orderBy, _ := set.GetString("order-by")
	formatString, _ := set.GetString("format")

	// Get the list of files
	files := internal.GitLsTree(repository, revision, extensions, languages)
	if len(exclude) != 0 {
		files = internal.Exclude(files, exclude)
	}
	if len(restrictTo) != 0 {
		files = internal.RestrictTo(files, restrictTo)
	}

	// Initialize statistics maps
	authors := make(map[string]*internal.Statistics)
	authorFiles := make(map[string]map[string]struct{}, 100)
	commitCounts := make(map[string]int)
	fileCounts := make(map[string]int)

	// Process each file and update statistics
	for _, file := range files {
		blameOutput := internal.GitBlame(repository, revision, file)
		authorData, commitData := internal.Blame(blameOutput, useCommitter, repository, revision, file)
		UpdateData(authorFiles, commitCounts, fileCounts, authorData, commitData)
	}

	internal.AuthorData(authors, authorFiles, commitCounts, fileCounts)

	sortedAuthors := internal.Sort(authors, orderBy)
	internal.Format(sortedAuthors, formatString)
}

func UpdateData(
	authorFiles map[string]map[string]struct{},
	commitCounts, fileCounts map[string]int,
	authors map[string][]string,
	commits map[string]int,
) {
	// Update commit counts
	for author, commitCount := range commits {
		commitCounts[author] += commitCount
	}

	// Update file counts and author files
	for author, files := range authors {
		fileCounts[author]++
		if _, exists := authorFiles[author]; !exists {
			authorFiles[author] = make(map[string]struct{})
		}
		for _, file := range files {
			authorFiles[author][file] = struct{}{}
		}
	}
}

func FlagSet() *flag.FlagSet {
	set := flag.NewFlagSet("flag_set", flag.ExitOnError)
	set.String("repository", ".", "path to Git repo")
	set.String("revision", "HEAD", "pointer on commit")
	set.String("order-by", "lines", "key for result sorting")
	set.Bool("use-committer", false, "flag to switch between author and committer")
	set.String("format", "tabular", "format of the output")
	set.StringSlice("extensions", []string{}, "list of extensions")
	set.StringSlice("languages", []string{}, "list of permitted languages")
	set.StringSlice("exclude", []string{}, "set of glob patterns ")
	set.StringSlice("restrict-to", []string{}, "set of glob patterns for restrict")

	err := set.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	return set
}
