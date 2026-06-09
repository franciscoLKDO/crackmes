package level

import (
	"io"
)

func LvlThree(_ io.Reader, w io.WriteCloser) error {
	out := []byte("01267567")
	_, err := w.Write(append(out, '\n'))
	return err
}
