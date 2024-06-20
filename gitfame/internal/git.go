package internal

import (
	"os/exec"
	"strings"
)

func GitLog(rep, commit, fileName string) (string, string) {
	cmd := exec.Command("git", "log", "--pretty=format:%H %an", commit, "--", fileName)
	cmd.Dir = rep
	output, _ := cmd.Output()

	s := strings.Split(string(output), "\n")[0]
	fields := strings.Split(s, " ")
	return fields[0], strings.Join(fields[1:], " ")
}

func GitBlame(repo, commit, fileName string) []string {
	cmd := exec.Command("git", "blame", "--porcelain", commit, fileName)
	cmd.Dir = repo
	blameOutput, _ := cmd.Output()

	return strings.FieldsFunc(string(blameOutput), func(r rune) bool {
		return r == '\n'
	})
}

func GitLsTree(repo, commit string, extensions, languages []string) []string {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", commit)
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var files []string
	for _, name := range strings.Split(string(out), "\n") {
		if name == "" {
			continue
		}
		files = append(files, name)
	}

	if len(languages) != 0 {
		getLanguagesExtensions(languages, &extensions)
	}

	if len(extensions) != 0 {
		files = filterExtension(files, extensions)
	}

	return files
}
