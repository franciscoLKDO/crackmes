package level

import (
	"io"
)

func LvlOne(_ io.Reader, w io.WriteCloser) error {
	in := "this is sparta"
	w.Write(append([]byte(in), '\n'))

	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[i] ^ 4
	}
	_, err := w.Write(append(out, '\n'))
	return err
}
