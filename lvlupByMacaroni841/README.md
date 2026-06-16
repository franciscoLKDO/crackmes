# LevelMacaroni

A Go-based challenge harness for the `lvlupByMacaroni841` binary.

The repository launches an external binary under a pseudo-terminal, inspects its output, and automatically responds to four challenge levels using dedicated handlers.

## Repository structure

- `main.go` — entry point; starts the target binary and routes output to level handlers.
- `engine/engine.go` — PTY engine that reads process output and dispatches request handlers based on level prompts.
- `level/lvl1.go` — Level 1 handler: reads a line, applies a fixed XOR, and writes the response.
- `level/lvl2.go` — Level 2 handler: generates a decoded response from a hardcoded obfuscated string.
- `level/lvl3.go` — Level 3 handler: returns a deterministic fixed code.
- `level/lvl4.go` — Level 4 handler: embeds a shared object hook, intercepts `time()`, and generates a key from the communicated epoch.
- `level/c/time.c` — C hook library used with `LD_PRELOAD` to intercept time calls from the target binary.

## How it works

1. `main.go` launches `./lvlupByMacaroni841` with an optional level argument.
2. It creates a PTY and uses `engine.New` to manage the binary's I/O.
3. When the binary prints one of the level prompts, the engine dispatches the corresponding handler.
4. `LvlFour` additionally starts a Unix socket server and injects `level/c/time.c` as a preloaded library to capture the target binary's `time()` value.

## Requirements

- Go `1.25` or later
- `gcc` for building the shared hook library (optional, embeded in go binary)
- `./lvlupByMacaroni841` downloaded from [crackmes.one](https://crackmes.one/crackme/66c724b9b899a3b9dd02ad98) 

## Build and run

1. Build the hook library:

   ```sh
   gcc -shared -fPIC -o level/c/hook.so level/c/time.c -ldl
   ```

2. Run the harness:

   ```sh
   go run main.go --lvlarg=2
   ```

   Use `--lvlarg=3` to pass a different startup argument to the target binary.

## Notes

- The target binary is not included in this repository. The harness expects `./lvlupByMacaroni841` to exist and be executable.
- Level 4 relies on the `LD_PRELOAD` hook and a Unix domain socket at `/tmp/keygen.sock` to communicate the current epoch.

## Learnings
### asm
- optimised modulo expression in asm
for example:
```
call   11d0 <rand@plt>
movsxd rdx,eax
imul   rdx,rdx,0x66666667
shr    rdx,0x20
sar    edx,1
mov    ecx,eax
sar    ecx,0x1f
sub    edx,ecx
mov    DWORD PTR [rbp-0x358],edx
mov    ecx,DWORD PTR [rbp-0x358]
mov    edx,eax
shl    edx,0x2
add    edx,ecx
sub    eax,edx
mov    DWORD PTR [rbp-0x358],eax
```

imul with weird const and shr/sar
then eax - n*edx -> modulo of %5 in this case 

optimised modulo for 2^n

```
shr    edx,0x1e
add    eax,edx
and    eax,0x3
sub    eax,edx
```
this is modulo 4 for signed
-> x and 3 (mask de 3)

  10101110
  00000011
->00000010 -> 2

we just watch the x last bits
---
### gdb
- Find an entrypoint of a stripped binary with gdb and `info file`
- tricks on gdb to redirect i/o 
```bash
# note the output of terminal destination
tty 
tail -f /dev/null

#then in gdb terminal
tty /dev/pts/X # the result of tty 
```
- save to file some logs
```bash
#gdb
set logging file <YourFile>
set logging enabled on
.....
set logging enabled off
```
- save breakpoints 
```bash 
#gdb
save breakpoints <my_file>
```