# go-xargs

A simple implement of xargs, written in golang

### Note.
> There need to specify command like passing '-bin ls' or '-bin echo'. diff from 'xargs'

### Example:

```bash
$ go run main.go -h
Usage of C:\Users\wurui\AppData\Local\Temp\go-build714031134\b001\exe\main.exe:
  -P int
        max-procs, default 0, It mean no limit (default 1)
  -bin string
        command to exec, default echo (default "echo")
  -n int
        max-args, default 0, It mean no limit (default 1)
exit status 2
```

```bash
$ echo {1..10} | go run main.go -P 3 -n 3
7 8 9
4 5 6
1 2 3
10
```

```bash
$ echo {1..10} | go run main.go -P 3 -n 2
9 10
3 4
1 2
5 67 8
```

### Time consuming:

```
2~3h
```

### Bug:

```bash 
1. It always echo an empty
```