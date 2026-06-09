package level

import (
	"io"
)

func LvlTwo(_ io.Reader, w io.WriteCloser) error {
	res := []byte("1d86ce")
	out := make([]byte, len(res))
	for i := range out {
		out[i] = (res[i] - byte(i)) ^ 0x12
	}

	_, err := w.Write(append(out, '\n'))
	return err
}
