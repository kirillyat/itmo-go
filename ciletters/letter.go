//go:build !solution

package ciletters

import (
	"strings"
	"text/template"
)

const notificationTemplate = `Your pipeline #{{ .Pipeline.ID }} {{if ne .Pipeline.Status "ok"}}has failed{{else}}passed{{end}}!
    Project:      {{ .Project.GroupID }}/{{.Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ slice .Commit.Hash 0 8 }} {{ .Commit.Message }}
    CommitAuthor: {{ .Commit.Author }}{{range .Pipeline.FailedJobs}}
        Stage: {{ .Stage }}, Job {{ .Name }}{{range splitRunnerLog .RunnerLog}}
            {{ . }}{{end}}
{{end}}`

func splitRunnerLog(logContent string) []string {
	logLines := strings.Split(logContent, "\n")
	if len(logLines) > 9 {
		logLines = logLines[9:]
	}
	return logLines
}

func MakeLetter(notification *Notification) (string, error) {
	templateFuncMap := template.FuncMap{
		"splitRunnerLog": splitRunnerLog,
	}

	tpl, err := template.New("notificationEmail").Funcs(templateFuncMap).Parse(notificationTemplate)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	if err = tpl.Execute(&builder, notification); err != nil {
		return "", err
	}

	return builder.String(), nil
}
