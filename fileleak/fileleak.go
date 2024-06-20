//go:build !solution

package fileleak

import (
	"io/fs"
	"os"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getOpenFiles() []fs.FileInfo {
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		panic(err)
	}

	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		var info fs.FileInfo
		info, err = entry.Info()

		if err != nil {
			continue
		}

		infos = append(infos, info)
	}

	return infos
}

func VerifyNone(t testingT) {
	before := getOpenFiles()
	t.Cleanup(func() {
		after := getOpenFiles()
		if !reflect.DeepEqual(after, before) {
			t.Errorf("Leaks detected. Before: %v, After: %v", before, after)
		}
	})
}
