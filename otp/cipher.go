//go:build !solution

package otp

import (
	"bytes"
	"io"
)

type cipherReader struct {
	r    io.Reader
	prng io.Reader
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &cipherReader{
		r:    r,
		prng: prng,
	}
}

func (c *cipherReader) Read(p []byte) (n int, err error) {
	n, errRead := c.r.Read(p)
	if errRead != nil && errRead != io.EOF {
		return 0, errRead
	}

	key := make([]byte, n)
	_, errPRNG := io.ReadFull(c.prng, key)
	if errPRNG != nil {
		return 0, errPRNG
	}

	for i := 0; i < n; i++ {
		p[i] ^= key[i]
	}

	return n, errRead
}

type cipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &cipherWriter{
		w:    w,
		prng: prng,
	}
}

func (c *cipherWriter) Write(p []byte) (int, error) {
	randomBytes, err := io.ReadAll(c.prng)
	if err != nil {
		return 0, err
	}

	var encrypted bytes.Buffer
	for i := range p {
		encrypted.WriteByte(p[i] ^ randomBytes[i])
	}

	return c.w.Write(encrypted.Bytes())
}
