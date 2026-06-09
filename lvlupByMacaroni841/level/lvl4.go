package level

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"lvlup/engine"
	"net"
	"os"
)

//go:embed c/hook.so
var hook []byte

func LvlFour(ctx context.Context) (string, engine.HandlerFunc, error) {
	hookPath, err := loadHook(ctx)
	if err != nil {
		return "", nil, err
	}
	chw := timeCatcher(ctx)
	return hookPath, func(_ io.Reader, w io.WriteCloser) error {
		_, err := w.Write((append(keygen(<-chw), '\n')))
		return err
	}, nil
}

func keygen(epoch int) []byte {
	mb := magicBytes()

	C.srand(C.uint(epoch))
	firstRand := C.rand() % 4
	secondRand := C.rand() % 3
	thirdRand := C.rand() % 5

	// I am hungry
	nutriscore := 0
	for i := 0; i <= 5; i++ {
		if (C.rand() % 40) < 0x15 {
			nutriscore++
		} else {
			nutriscore--
		}
	}

	secret := make([]byte, 6)
	secret[0] = mb[2]
	secret[4] = mb[6]
	if nutriscore < 0 {
		secret[0] = mb[3]
		secret[4] = mb[1]
	}

	switch secondRand {
	case 0:
		secret[5] = mb[4]
	case 1:
		secret[5] = mb[1]
	case 2:
		secret[5] = mb[0]
	default:
		fmt.Println("error out of track inner 1")
		return []byte{}
	}

	a, b, c, d := 2, 4, 9, 9
	secret[2] = mb[2]
	if thirdRand <= 2 {
		secret[2] = mb[8]
		a, b, c, d = 2, 7, 4, 8
	}

	secret[1] = mb[c]
	secret[3] = mb[d]
	if firstRand <= 1 {
		secret[1] = mb[a]
		secret[3] = mb[b]
	}

	return secret
}

func magicBytes() []byte {
	mn := make([]byte, 16)
	binary.LittleEndian.PutUint64(mn, 0x4a2a42515a2b2440)
	binary.LittleEndian.PutUint64(mn[8:], 0x3e3c3f3a7b7c2526)
	return mn
}

func loadHook(ctx context.Context) (string, error) {
	file, err := os.CreateTemp("/tmp", "")
	if err != nil {
		return "", err
	}
	if _, err := file.Write(hook); err != nil {
		return "", err
	}
	go func() {
		<-ctx.Done()
		os.Remove(file.Name())
	}()
	return file.Name(), nil
}

const socket = "/tmp/keygen.sock"

func timeCatcher(ctx context.Context) chan int {
	c := make(chan int)
	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Fatal(err)
			}

			defer conn.Close()
			buf := make([]byte, 8)

			for {
				if err := binary.Read(conn, binary.LittleEndian, &buf); err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(errors.Join(errors.New("error reading socket"), err))
				}
			}
			select {
			case <-ctx.Done():
				l.Close()
				close(c)
				os.Remove(socket)
				return
			case c <- int(binary.LittleEndian.Uint64(buf)):
			}
		}
	}()
	return c
}
