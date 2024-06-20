//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type bufferedLineReader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) LineReader {
	return &bufferedLineReader{
		reader: bufio.NewReader(r),
	}
}

func (blr *bufferedLineReader) ReadLine() (string, error) {
	var sb strings.Builder

	for {
		b, err := blr.reader.ReadByte()
		if err != nil {
			return sb.String(), err
		}
		if b == '\n' {
			break
		}
		sb.WriteByte(b)
	}

	return sb.String(), nil
}

type textLineWriter struct {
	writer io.Writer
}

func NewWriter(w io.Writer) LineWriter {
	return &textLineWriter{
		writer: w,
	}
}

func (tlw *textLineWriter) Write(line string) error {
	_, err := tlw.writer.Write([]byte(line + "\n"))
	return err
}

type heapItem struct {
	reader LineReader
	line   string
}

type minHeap []heapItem

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].line < h[j].line }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(heapItem))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &minHeap{}
	for _, r := range readers {
		str, err := r.ReadLine()
		if err != nil && (err != io.EOF) || (err == io.EOF && str == "") {
			continue
		}
		heap.Push(h, heapItem{
			line:   str,
			reader: r,
		})
	}
	heap.Init(h)

	for h.Len() > 0 {
		minheapItem := heap.Pop(h).(heapItem)

		err := w.Write(minheapItem.line)

		if err != nil {
			return err
		}

		str, err := minheapItem.reader.ReadLine()
		if err == nil || (str != "" && err == io.EOF) {
			heap.Push(h, heapItem{
				line:   str,
				reader: minheapItem.reader,
			})
			heap.Fix(h, h.Len()-1)
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	var readers []LineReader
	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		lr := NewReader(f)
		var lines []string
		for {
			str, errRL := lr.ReadLine()

			if errRL == io.EOF && str != "" {
				lines = append(lines, str)
			}

			if errRL != nil {
				break
			}

			lines = append(lines, str)
		}

		sort.Strings(lines)
		f.Close()
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)

		if err != nil {
			return err
		}

		lw := NewWriter(f)

		for _, str := range lines {
			err := lw.Write(str)
			if err != nil {
				return err
			}
		}
		f.Close()
	}

	lw := NewWriter(w)

	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		lr := NewReader(f)
		readers = append(readers, lr)
	}

	return Merge(lw, readers...)
}
