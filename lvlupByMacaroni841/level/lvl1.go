package level

import (
	"bufio"
	"io"
)

func LvlOne(r io.Reader, w io.WriteCloser) error {
	in, _, err := bufio.NewReader(r).ReadLine()
	if err != nil {
		panic(err)
	}

	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[i] ^ 4
	}
	_, err = w.Write(append(out, '\n'))
	return err
}
