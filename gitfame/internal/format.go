package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type AuthorJSON struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func Format(a AuthorSlice, format string) {
	switch format {
	case "tabular":
		Tabular(a)
	case "csv":
		CSV(a)
	case "json":
		JSON(a)
	case "json-lines":
		JSONLines(a)
	default:
		panic("format flag value is incorrect")
	}
}

func Tabular(a AuthorSlice) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")

	for _, i := range a.Slice {
		s := ""
		s += i.Name + "\t"
		s += strconv.Itoa(i.Statistics.Lines) + "\t"
		s += strconv.Itoa(i.Statistics.Commits) + "\t"
		s += strconv.Itoa(i.Statistics.Files)
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

func CSV(a AuthorSlice) {
	w := csv.NewWriter(os.Stdout)
	err := w.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		panic(err)
	}

	for _, i := range a.Slice {
		err = w.Write([]string{i.Name,
			strconv.Itoa(i.Statistics.Lines),
			strconv.Itoa(i.Statistics.Commits),
			strconv.Itoa(i.Statistics.Files)},
		)
		if err != nil {
			panic(err)
		}
	}

	w.Flush()
}

func JSON(a AuthorSlice) {
	authorJSONList := make([]AuthorJSON, len(a.Slice))
	for j, i := range a.Slice {
		authorJSONList[j] = AuthorJSON{
			Name:    i.Name,
			Lines:   i.Statistics.Lines,
			Commits: i.Statistics.Commits,
			Files:   i.Statistics.Files,
		}
	}
	b, _ := json.Marshal(&authorJSONList)
	os.Stdout.Write(b)

}

func JSONLines(a AuthorSlice) {
	for _, i := range a.Slice {
		b, _ := json.Marshal(AuthorJSON{
			Name:    i.Name,
			Lines:   i.Statistics.Lines,
			Commits: i.Statistics.Commits,
			Files:   i.Statistics.Files,
		})
		os.Stdout.Write(b)
		fmt.Println()
	}
}
