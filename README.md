# go-xargs

A simple implement of xargs, written in golang

### Note.
> There needs to specify command like passing '-bin ls' or '-bin echo'. This's different from 'xargs'

### Example:

```bash
$ go run main.go -h
Usage of /var/folders/y3/q1jmb_7s3jsfjh0t0y5xnd0c0000gn/T/go-build080969720/b001/exe/main:
  -P int
        max-procs, default 1 (default 1)
  -bin string
        command to exec, default echo (default "echo")
  -n int
        max-args, default 1  (default 1)
exit status 2
```

```bash
$ echo {1..10} | go run main.go -P 3 -n 3
7 8 9
4 5 6
1 2 3
10
```
