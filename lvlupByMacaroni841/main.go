package main

import (
	"context"
	"flag"
	"lvlup/engine"
	"lvlup/level"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const (
	appPath = "./lvlupByMacaroni841"
	hook    = "./c/hook.so"
)

func main() {
	lvlArg := flag.Int("lvlarg", 2, "set to 2 to run lvl2 or 3 for lvl3 (need to set binary arg to work)")
	flag.Parse()
	arg := "e"
	if lvlArg != nil && *lvlArg == 3 {
		arg = "0x3"
	}

	start(arg)
}

func start(arg string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := exec.Command(appPath, arg)
	env := os.Environ()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	unixSocket, lvl4, err := level.LvlFour(ctx)
	if err != nil {
		panic(err)
	}
	c.Env = append(env, "LD_PRELOAD="+unixSocket)

	eng, err := engine.New(c)
	if err != nil {
		panic(err)
	}

	eng.Handle("Level 1 -", level.LvlOne)
	eng.Handle("Level 2 -", level.LvlTwo)
	eng.Handle("Level 3 -", level.LvlThree)
	eng.Handle("Level 4 -", lvl4)

	// Handle pty size.
	_ = eng.Serve()
	time.Sleep(100 * time.Millisecond)
}
