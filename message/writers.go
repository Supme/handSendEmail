package message

import (
	"encoding/base64"
	"io"
	"os"
	"strings"
)

type delimitWriter struct {
	n      int
	cnt    int
	dr     []byte
	writer io.Writer
}

func newDelimitWriter(writer io.Writer, dr []byte, cnt int) *delimitWriter {
	return &delimitWriter{n: 0, cnt: cnt, dr: dr, writer: writer}
}

func (w *delimitWriter) Write(p []byte) (n int, err error) {
	for i := range p {
		_, err = w.writer.Write(p[i : i+1])
		if err != nil {
			break
		}
		if w.n++; w.n%w.cnt == 0 {
			w.writer.Write(w.dr)

		}
	}
	return w.n, err
}

func base64FileWriter(w io.Writer, f *os.File) (err error) {
	dwr := newDelimitWriter(w, []byte{0x0d, 0x0a}, 76) // 76 from RFC
	b64Enc := base64.NewEncoder(base64.StdEncoding, dwr)
	_, err = io.Copy(b64Enc, f)
	if err != nil {
		return err
	}

	return b64Enc.Close()
}

func base64TextWriter(w io.Writer, text string) (err error) {
	dwr := newDelimitWriter(w, []byte{0x0d, 0x0a}, 76) // 76 from RFC
	b64Enc := base64.NewEncoder(base64.StdEncoding, dwr)
	reader := strings.NewReader(text)
	_, err = io.Copy(b64Enc, reader)
	if err != nil {
		return err
	}

	return b64Enc.Close()
}
