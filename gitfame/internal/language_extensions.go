package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Language struct {
	Name       string
	Type       string
	Extensions []string
}

func getLanguagesExtensions(languages []string, extensions *[]string) {
	var l []Language
	data, err := os.ReadFile("../../configs/language_extensions.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &l); err != nil {
		panic(err)
	}

	set := make(map[string]int, len(l))

	for i, j := range l {
		set[strings.ToLower(j.Name)] = i
	}

	for _, i := range languages {
		if itr, ok := set[strings.ToLower(i)]; ok {
			*extensions = append(*extensions, l[itr].Extensions...)
		}
	}
}

func filterExtension(files, extensions []string) []string {
	extSet := make(map[string]struct{}, len(extensions))
	for _, ext := range extensions {
		extSet[ext] = struct{}{}
	}

	var result []string
	for _, file := range files {
		if _, exists := extSet[filepath.Ext(file)]; exists {
			result = append(result, file)
		}
	}

	return result
}
