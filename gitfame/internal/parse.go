package internal

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Statistics struct {
	Lines   int
	Commits int
	Files   int
}

type Author struct {
	Statistics Statistics
	Name       string
}

type AuthorSlice struct {
	Slice   []Author
	orderBy string
}

func (a AuthorSlice) Len() int {
	return len(a.Slice)
}

func (a AuthorSlice) Swap(i, j int) {
	a.Slice[i], a.Slice[j] = a.Slice[j], a.Slice[i]
}

func (a AuthorSlice) Less(i, j int) bool {
	var key1, key2 []int

	switch a.orderBy {
	case "lines":
		key1 = []int{a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Commits, a.Slice[i].Statistics.Files}
		key2 = []int{a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Commits, a.Slice[j].Statistics.Files}
	case "commits":
		key1 = []int{a.Slice[i].Statistics.Commits, a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Files}
		key2 = []int{a.Slice[j].Statistics.Commits, a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Files}
	case "files":
		key1 = []int{a.Slice[i].Statistics.Files, a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Commits}
		key2 = []int{a.Slice[j].Statistics.Files, a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Commits}
	default:
		panic("order-by flag value is incorrect")
	}

	for idx := range key1 {
		if key1[idx] != key2[idx] {
			return key1[idx] > key2[idx]
		}
	}

	return strings.ToLower(a.Slice[i].Name) < strings.ToLower(a.Slice[j].Name)
}

func RestrictTo(files, patterns []string) []string {
	var result []string
	for _, file := range files {
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, file); matched {
				result = append(result, file)
				break
			}
		}
	}
	return result
}

func Exclude(files, patterns []string) []string {
	var result []string
	for _, file := range files {
		exclude := false
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, file); matched {
				exclude = true
				break
			}
		}
		if !exclude {
			result = append(result, file)
		}
	}
	return result
}

func Blame(out []string, useCommitter bool, repo, commit, fileName string) (map[string][]string, map[string]int) {
	authors := make(map[string][]string)
	commits := make(map[string]int)

	if len(out) == 0 {
		hash, author := GitLog(repo, commit, fileName)
		commits[hash] = 0
		authors[author] = append(authors[author], hash)
		return authors, commits
	}

	isNextHash := true
	itr := 0
	var isWaitForAuthor bool
	var lastHash string

	for _, line := range out {
		if isNextHash {
			isNextHash = false
			if itr == 0 {
				parts := strings.Split(line, " ")
				itr, _ = strconv.Atoi(parts[len(parts)-1])
				lastHash = parts[0]
				commits[lastHash] += itr
				isWaitForAuthor = true
			}
			itr--
		} else if line[0] != '\t' && isWaitForAuthor {
			parts := strings.Split(line, " ")
			var prefix string
			if useCommitter {
				prefix = "committer"
			} else {
				prefix = "author"
			}
			if parts[0] == prefix {
				name := line[len(prefix)+1:]
				authors[name] = append(authors[name], lastHash)
				isWaitForAuthor = false
			}
		} else if line[0] == '\t' {
			isNextHash = true
		}
	}

	return authors, commits
}

func AuthorData(authors map[string]*Statistics, a map[string]map[string]struct{}, c, files map[string]int) {
	for i, j := range a {
		authors[i] = &Statistics{
			Lines:   0,
			Commits: len(j),
			Files:   files[i],
		}
	}
	for i, j := range a {
		lines := 0
		for k := range j {
			lines += c[k]
		}
		authors[i].Lines += lines
	}
}

func Sort(authors map[string]*Statistics, orderBy string) AuthorSlice {
	var authorSlice AuthorSlice
	authorSlice.orderBy = orderBy
	for i, j := range authors {
		authorSlice.Slice = append(authorSlice.Slice, Author{
			Statistics: *j,
			Name:       i,
		})
	}

	sort.Sort(authorSlice)
	return authorSlice
}
