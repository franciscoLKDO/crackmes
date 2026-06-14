package main

import (
	"context"
	"flag"
	"fmt"
	"lvlup/engine"
	"lvlup/level"
	"os"
	"os/exec"
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

	if err := start(arg); err != nil {
		panic(err)
	}
}

func start(arg string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := exec.Command(appPath, arg)
	env := os.Environ()

	lvl4, err := level.NewLvlFour(ctx)
	if err != nil {
		return err
	}
	c.Env = append(env, "LD_PRELOAD="+lvl4.HookFile)

	session, err := engine.NewSession(c, os.Stdout)
	if err != nil {
		return err
	}
	eng := engine.NewEngine(&session)

	eng.Register("Level 1 -", level.LvlOne)
	eng.Register("Level 2 -", level.LvlTwo)
	eng.Register("Level 3 -", level.LvlThree)
	eng.Register("Level 4 -", lvl4.Handler)

	go eng.Start()

	if err = session.Wait(); err != nil {
		fmt.Println("process exited with error:", err)
	}

	lvl4.Clean()
	return session.Close()
}
